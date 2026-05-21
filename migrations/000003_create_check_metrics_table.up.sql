CREATE TABLE check_metrics (
    check_id bigint PRIMARY KEY REFERENCES checks(id) ON DELETE CASCADE,
    line_count bigint NOT NULL,
    comment_line_count bigint NOT NULL,
    comment_ratio double precision NOT NULL,
    function_count bigint NOT NULL,
    average_function_length double precision NOT NULL,
    max_function_length bigint NOT NULL,
    conditional_count bigint NOT NULL,
    loop_count bigint NOT NULL,
    max_nesting_depth bigint NOT NULL,
    global_variable_count bigint NOT NULL,
    long_line_count bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);