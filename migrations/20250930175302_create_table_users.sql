-- +goose Up
-- +goose StatementBegin
create table url_shorten(
    url         text      not null,
    alias       text      not null,
    created_at  timestamp not null,
    updated_at  timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
