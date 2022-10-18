-- +goose Up
-- +goose StatementBegin
create table events
(
    id                    uuid         not null
        constraint events_pk
            primary key,
    title                 varchar(255) not null,
    description           text,
    start_at              timestamp    not null,
    end_at                timestamp      not null,
    user_id               uuid         not null,
    notification_duration bigint
);

alter table events
    owner to calendar;

create index events_user_id_start_at_end_at_index
    on events (user_id, start_at, end_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
