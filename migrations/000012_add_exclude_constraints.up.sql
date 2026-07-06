CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Запрет пересечения шаблонов работы по одному мастеру
-- (в одном weekday, с пересекающимся периодом действия и пересекающимися часами смены)
ALTER TABLE "hairdresser_work_patterns"
    ADD CONSTRAINT "excl_pattern_overlap"
    EXCLUDE USING gist (
        hairdresser_id WITH =,
        weekday WITH =,
        daterange(effective_from, COALESCE(effective_to, 'infinity'::date), '[]') WITH &&,
        tsrange(
            ('2000-01-01'::date + shift_start)::timestamp,
            ('2000-01-01'::date + shift_end)::timestamp
        ) WITH &&
    );

-- Один мастер не может быть в двух местах одновременно
ALTER TABLE "hairdresser_schedules"
    ADD CONSTRAINT "excl_schedule_overlap"
    EXCLUDE USING gist (
        hairdresser_id WITH =,
        tstzrange(shift_start, shift_end) WITH &&
    );

-- Запрет двойных записей к одному мастеру (кроме отменённых/неявок)
ALTER TABLE "bookings"
    ADD CONSTRAINT "excl_booking_overlap"
    EXCLUDE USING gist (
        hairdresser_id WITH =,
        tstzrange(starts_at, ends_at) WITH &&
    ) WHERE (status NOT IN ('cancelled', 'no_show'));
