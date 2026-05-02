import { DynamoDBClient } from '@aws-sdk/client-dynamodb';
import { DynamoDBDocumentClient, GetCommand, ScanCommand } from '@aws-sdk/lib-dynamodb';

type Query = Record<string, string | undefined>;
type JsonBody = Record<string, unknown>;

type LambdaEvent = {
  requestContext?: { http?: { method?: string } };
  httpMethod?: string;
  rawPath?: string;
  path?: string;
  queryStringParameters?: Query;
};

const tableName = process.env.JOBS_TABLE_NAME;
const resultsBaseUrl = (process.env.RESULTS_BASE_URL || '').replace(/\/$/, '');

const ddb = DynamoDBDocumentClient.from(new DynamoDBClient({}));

const STATUS_MAP = new Map<string, string>([
  ['done', 'succeeded'],
  ['completed', 'succeeded'],
  ['success', 'succeeded'],
  ['succeeded', 'succeeded'],
  ['failed', 'failed'],
  ['error', 'failed'],
  ['cancelled', 'canceled'],
  ['canceled', 'canceled'],
  ['in-progress', 'running'],
  ['in_progress', 'running'],
  ['running', 'running'],
  ['processing', 'running'],
  ['queued', 'queued'],
  ['pending', 'queued'],
]);

function json(statusCode: number, body: JsonBody) {
  return {
    statusCode,
    headers: {
      'content-type': 'application/json; charset=utf-8',
      'access-control-allow-origin': '*',
      'access-control-allow-methods': 'GET,OPTIONS',
    },
    body: JSON.stringify(body),
  };
}

function normalizeStatus(raw: unknown) {
  const key = String(raw || '').trim().toLowerCase();
  return STATUS_MAP.get(key) || 'queued';
}

function inferProgress(status: string) {
  if (status === 'succeeded') return 100;
  if (status === 'running') return 50;
  if (status === 'queued') return 10;
  return 0;
}

function parseTime(raw: unknown) {
  const s = String(raw || '').trim();
  if (!s) return null;

  const asNumber = Number(s);
  if (Number.isFinite(asNumber) && String(Math.trunc(asNumber)) === s) {
    return new Date(asNumber * 1000).toISOString();
  }

  const d = new Date(s);
  if (!Number.isNaN(d.getTime())) return d.toISOString();
  return null;
}

function inferOutputFiles(row: Record<string, unknown>) {
  const prefix = String(row?.s3Output || '').trim();
  if (!prefix) return [];

  const clean = prefix.replace(/^s3:\/\//, '');
  const slashIndex = clean.indexOf('/');
  if (slashIndex < 0) return [];

  let base = clean.slice(slashIndex + 1).replace(/\/$/, '');
  const jobId = String(row?.uuid || '').trim();
  if (jobId && !base.includes(jobId)) {
    base = `${base}/${jobId}`;
  }

  return [
    { key: `${base}/model.splat`, file_name: 'model.splat' },
    { key: `${base}/model.ply`, file_name: 'model.ply' },
    { key: `${base}/model.spz`, file_name: 'model.spz' },
  ];
}

function toSummary(row: Record<string, unknown>) {
  const status = normalizeStatus(row?.uuidStatus);
  const updatedAt = parseTime(row?.updatedAt) || parseTime(row?.endTimestamp) || parseTime(row?.startTimestamp) || new Date().toISOString();
  const createdAt = parseTime(row?.startTimestamp) || updatedAt;

  return {
    job_id: String(row?.uuid || ''),
    status,
    progress_percent: inferProgress(status),
    created_at: createdAt,
    updated_at: updatedAt,
  };
}

function toDetails(row: Record<string, unknown>) {
  const summary = toSummary(row);
  const outputFiles = inferOutputFiles(row);

  return {
    summary,
    attempt: 1,
    source_ref: String(row?.s3Input || ''),
    error_message: row?.errorMsg ? String(row.errorMsg) : undefined,
    started_at: parseTime(row?.startTimestamp) || undefined,
    finished_at: parseTime(row?.endTimestamp) || undefined,
    output_files: outputFiles,
  };
}

function toResultUrls(details: ReturnType<typeof toDetails>, ttlSeconds: number) {
  const ttl = Math.max(60, Math.min(86400, ttlSeconds || 3600));
  const expiresAt = new Date(Date.now() + ttl * 1000).toISOString();

  return {
    items: (details.output_files || []).map((f) => ({
      key: f.key,
      file_name: f.file_name,
      url: `${resultsBaseUrl}/${f.key}`,
      expires_at: expiresAt,
    })),
  };
}

async function listJobs(query: Query) {
  const limit = Math.max(1, Math.min(200, Number(query?.limit || 100)));
  const offset = Math.max(0, Number(query?.offset || 0));
  const statusFilter = query?.status ? String(query.status).trim().toLowerCase() : '';

  const out = await ddb.send(new ScanCommand({ TableName: tableName }));
  let items = (out.Items || []).map((x) => toSummary(x as Record<string, unknown>));

  if (statusFilter) {
    items = items.filter((x) => x.status === statusFilter);
  }

  items.sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime());
  return { items: items.slice(offset, offset + limit) };
}

async function getJob(jobId: string) {
  const out = await ddb.send(new GetCommand({
    TableName: tableName,
    Key: { uuid: jobId },
  }));

  if (!out.Item) {
    return null;
  }

  return toDetails(out.Item as Record<string, unknown>);
}

export async function handler(event: LambdaEvent) {
  try {
    if (!tableName) {
      return json(500, { error: 'JOBS_TABLE_NAME is not set' });
    }

    const method = event?.requestContext?.http?.method || event?.httpMethod || 'GET';
    const path = event?.rawPath || event?.path || '/';

    if (method === 'OPTIONS') return json(200, { ok: true });
    if (method !== 'GET') return json(405, { error: 'method not allowed' });

    if (path === '/healthz') {
      return json(200, { status: 'ok' });
    }

    if (path === '/v1/jobs') {
      const payload = await listJobs(event?.queryStringParameters || {});
      return json(200, payload);
    }

    const resultUrlsMatch = path.match(/^\/v1\/jobs\/([^/]+)\/result-urls$/);
    if (resultUrlsMatch) {
      const jobId = decodeURIComponent(resultUrlsMatch[1]);
      const details = await getJob(jobId);
      if (!details) return json(404, { error: 'job not found' });
      const ttlSeconds = Number(event?.queryStringParameters?.ttl_seconds || 3600);
      return json(200, toResultUrls(details, ttlSeconds));
    }

    const detailsMatch = path.match(/^\/v1\/jobs\/([^/]+)$/);
    if (detailsMatch) {
      const jobId = decodeURIComponent(detailsMatch[1]);
      const details = await getJob(jobId);
      if (!details) return json(404, { error: 'job not found' });
      return json(200, details);
    }

    return json(404, { error: 'not found' });
  } catch (error) {
    console.error(error);
    return json(500, { error: 'internal server error' });
  }
}
