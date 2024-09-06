CREATE TABLE IF NOT EXISTS casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100),
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);

INSERT INTO casbin_rule (ptype, v0, v1, v2) 
VALUES ('p', 'admin', 'data1', 'read');
INSERT INTO casbin_rule (ptype, v0, v1) 
VALUES ('g', 'alice', 'admin');