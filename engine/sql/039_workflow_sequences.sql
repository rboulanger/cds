-- +migrate Up
CREATE TABLE IF NOT EXISTS "workflow_sequences" (
    workflow_id BIGINT,
    current_val BIGINT,
    PRIMARY KEY (workflow_id)
);

SELECT create_foreign_key_idx_cascade('FK_WORKFLOW_SEQUANCES_WORKFLOW', 'workflow_sequences', 'workflow', 'workflow_id', 'id');

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION workflow_sequences_nextval(w_id integer) RETURNS integer AS $$
DECLARE
    workflow_exists integer;
    cur_val integer;
BEGIN
    SELECT    count(1) INTO WORKFLOW_EXISTS
    FROM      workflow_sequences
    WHERE     workflow_id = $1;

    IF WORKFLOW_EXISTS = 0 THEN
        BEGIN
            INSERT INTO workflow_sequences(workflow_id, current_val) VALUES ($1, 0);
        EXCEPTION WHEN others THEN
            -- Do nothing
        END;
    END IF;
    
    SELECT    current_val INTO CUR_VAL
    FROM      workflow_sequences
    WHERE     workflow_id = $1 FOR UPDATE;

    UPDATE    workflow_sequences SET current_val = CUR_VAL + 1;

    RETURN    CUR_VAL + 1;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE "workflow_sequences";
DROP FUNCTION workflow_sequences_nextval(integer);