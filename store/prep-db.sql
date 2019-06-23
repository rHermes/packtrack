CREATE TABLE IF NOT EXISTS trackers (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT NOT NULL,
	url TEXT NOT NULL
);


-- TODO(rHermes): Asses if we should have a work_
CREATE TABLE IF NOT EXISTS work_items (
	id BIGSERIAL UNIQUE NOT NULL,
	created TIMESTAMPTZ NOT NULL,
	tracker INTEGER NOT NULL REFERENCES trackers(id),
	item TEXT NOT NULL,
	PRIMARY KEY (created, tracker, item
);


CREATE TABLE IF NOT EXISTS work_status (
	node TEXT,
	work_item BIGINT NOT NULL REFERENCSE work_items(id),

	
)

CREATE TABLE IF NOT EXISTS raws (
	ts TIMESTAMPTZ,
	node TEXT,
	worker TEXT,
	work_item BIGINT NOT NULL REFERENCES work_items(id),
	js jsonb NOT NULL,
	PRIMARY KEY (ts, node, worker, work_item)
);
