INSERT INTO apps (id, name, secret)
VALUES ('930b867c-97b5-48e9-91f2-55668aa47d41', 'test', 'test-secret')
    ON CONFLICT DO NOTHING;