ALTER TABLE "bookings" DROP CONSTRAINT IF EXISTS "excl_booking_overlap";
ALTER TABLE "hairdresser_schedules" DROP CONSTRAINT IF EXISTS "excl_schedule_overlap";
ALTER TABLE "hairdresser_work_patterns" DROP CONSTRAINT IF EXISTS "excl_pattern_overlap";

DROP EXTENSION IF EXISTS btree_gist;
