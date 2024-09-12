CREATE TABLE "books" (
    "id" BIGSERIAL,
    "created_at" TIMESTAMPTZ,
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    "thumbnail" TEXT,
    "detail_url" TEXT,
    "title" TEXT,
    "rating" BIGINT,
    "price" TEXT,
    "instock" BOOLEAN,
    PRIMARY KEY ("id"),
    CONSTRAINT "uni_books_detail_url" UNIQUE ("detail_url")
);

CREATE INDEX IF NOT EXISTS "idx_books_deleted_at" ON "books" ("deleted_at")