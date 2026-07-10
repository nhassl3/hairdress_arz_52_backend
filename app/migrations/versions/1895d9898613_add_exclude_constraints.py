"""add_exclude_constraints

Revision ID: 1895d9898613
Revises: cba5c918b793
Create Date: 2026-04-21 23:55:49.079241

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = '1895d9898613'
down_revision: Union[str, Sequence[str], None] = 'cba5c918b793'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
     # Включить расширение btree_gist
    op.execute("CREATE EXTENSION IF NOT EXISTS btree_gist")

    # Запрет пересечения шаблонов работы по одному мастеру
    op.execute("""
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
        )
    """)

    # Один мастер не может быть в двух местах одновременно
    op.execute("""
        ALTER TABLE "hairdresser_schedules"
        ADD CONSTRAINT "excl_schedule_overlap"
        EXCLUDE USING gist (
            hairdresser_id WITH =,
            tstzrange(shift_start, shift_end) WITH &&
        )
    """)

    # Запрет двойных записей к одному мастеру (кроме отменённых/неявок)
    op.execute("""
        ALTER TABLE "bookings"
        ADD CONSTRAINT "excl_booking_overlap"
        EXCLUDE USING gist (
            hairdresser_id WITH =,
            tstzrange(starts_at, ends_at) WITH &&
        )
        WHERE (status NOT IN ('cancelled', 'no_show'))
    """)


def downgrade() -> None:
    op.execute("ALTER TABLE bookings DROP CONSTRAINT IF EXISTS excl_booking_overlap")
    op.execute("ALTER TABLE hairdresser_schedules DROP CONSTRAINT IF EXISTS excl_schedule_overlap")
    op.execute("ALTER TABLE hairdresser_work_patterns DROP CONSTRAINT IF EXISTS excl_pattern_overlap")
    op.execute("DROP EXTENSION IF EXISTS btree_gist")
