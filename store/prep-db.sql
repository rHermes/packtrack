CREATE TABLE IF NOT EXISTS trackers (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT NOT NULL,
	url TEXT NOT NULL
);

INSERT INTO 
	trackers (name, description, url)
VALUES
	('bring', 'the norwegian postal service', 'https://developer.bring.com/')
ON CONFLICT
	DO NOTHING;

CREATE TABLE IF NOT EXISTS scrape_jobs (
	id BIGSERIAL PRIMARY KEY,
	tracker INTEGER NOT NULL REFERENCES trackers(id),
	args JSONB DEFAULT '{}',
	status TEXT NOT NULL DEFAULT 'created',
	created_at TIMESTAMPTZ NOT NULL,
	start_time TIMESTAMPTZ,
	end_time TIMESTAMPTZ,
	stats JSONB DEFAULT '{}',
	resp JSONB,
	CONSTRAINT created_before_started CHECK (created_at <= start_time),
	CONSTRAINT started_before_ended CHECK (start_time <= end_time),
	CONSTRAINT end_must_start CHECK ( (end_time IS NULL) OR (start_time IS NOT NULL))
);
CREATE INDEX IF NOT EXISTS idx_scrape_jobs_status ON scrape_jobs (status);
CREATE INDEX IF NOT EXISTS idx_scrape_jobs_created_at ON scrape_jobs (created_at);
CREATE INDEX IF NOT EXISTS idx_scrape_jobs_end_time ON scrape_jobs (end_time);
CREATE INDEX IF NOT EXISTS idx_scrape_jobs_start_time ON scrape_jobs (start_time);
