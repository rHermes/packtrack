-- Copyright (c) 2019 Teodor Spæren
--
-- Permission is hereby granted, free of charge, to any person obtaining a copy of
-- this software and associated documentation files (the "Software"), to deal in
-- the Software without restriction, including without limitation the rights to
-- use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
-- the Software, and to permit persons to whom the Software is furnished to do so,
-- subject to the following conditions:
--
-- The above copyright notice and this permission notice shall be included in all
-- copies or substantial portions of the Software.
--
-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
-- IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
-- FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
-- COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
-- IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
-- CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
CREATE INDEX IF NOT EXISTS idx_scrape_jobs_id_where_status_eq_created ON scrape_jobs(id) WHERE status = 'created';
