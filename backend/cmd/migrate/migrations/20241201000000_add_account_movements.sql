-- Create "account_movements" table
CREATE TABLE "public"."account_movements" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "account_id" uuid NOT NULL,
  "amount" numeric(15,2) NOT NULL,
  "occurred_at" timestamptz NOT NULL,
  "note" character varying(255) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_account_movements_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_account_movements_account" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_account_movements_user_id" to table: "account_movements"
CREATE INDEX "idx_account_movements_user_id" ON "public"."account_movements" ("user_id");
-- Create index "idx_account_movements_account_id" to table: "account_movements"
CREATE INDEX "idx_account_movements_account_id" ON "public"."account_movements" ("account_id");
-- Create index "idx_account_movements_occurred_at" to table: "account_movements"
CREATE INDEX "idx_account_movements_occurred_at" ON "public"."account_movements" ("occurred_at");

-- Create updated_at trigger for account_movements table
CREATE TRIGGER update_account_movements_updated_at BEFORE UPDATE ON "public"."account_movements" FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();