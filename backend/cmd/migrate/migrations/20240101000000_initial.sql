-- Create trigger function for updating timestamps
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "auth0_id" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_auth0_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_auth0_id" ON "public"."users" ("auth0_id");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");

-- Create "category_masters" table
CREATE TABLE "public"."category_masters" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "name" character varying(100) NOT NULL,
  "type" character varying(20) NOT NULL,
  "icon" character varying(50) NULL,
  "color" character varying(7) NULL,
  "display_order" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
-- Create index "idx_category_masters_type" to table: "category_masters"
CREATE INDEX "idx_category_masters_type" ON "public"."category_masters" ("type");

-- Create "accounts" table
CREATE TABLE "public"."accounts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "type" character varying(50) NOT NULL,
  "balance" numeric(15,2) NOT NULL DEFAULT 0.00,
  "currency" character varying(3) NOT NULL DEFAULT 'JPY',
  "is_active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_accounts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_accounts_user_id" to table: "accounts"
CREATE INDEX "idx_accounts_user_id" ON "public"."accounts" ("user_id");

-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "category_master_id" uuid NOT NULL,
  "custom_name" character varying(100) NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_master" FOREIGN KEY ("category_master_id") REFERENCES "public"."category_masters" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_categories_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_categories_master_id" to table: "categories"
CREATE INDEX "idx_categories_master_id" ON "public"."categories" ("category_master_id");
-- Create index "idx_categories_user_id" to table: "categories"
CREATE INDEX "idx_categories_user_id" ON "public"."categories" ("user_id");
-- Create index "idx_categories_user_master" to table: "categories"
CREATE UNIQUE INDEX "idx_categories_user_master" ON "public"."categories" ("user_id", "category_master_id");

-- Create "transactions" table
CREATE TABLE "public"."transactions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "account_id" uuid NOT NULL,
  "category_id" uuid NOT NULL,
  "amount" numeric(15,2) NOT NULL,
  "type" character varying(20) NOT NULL,
  "date" date NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_transactions_account" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_transactions_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_transactions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_transactions_account_id" to table: "transactions"
CREATE INDEX "idx_transactions_account_id" ON "public"."transactions" ("account_id");
-- Create index "idx_transactions_category_id" to table: "transactions"
CREATE INDEX "idx_transactions_category_id" ON "public"."transactions" ("category_id");
-- Create index "idx_transactions_date" to table: "transactions"
CREATE INDEX "idx_transactions_date" ON "public"."transactions" ("date");
-- Create index "idx_transactions_user_date" to table: "transactions"
CREATE INDEX "idx_transactions_user_date" ON "public"."transactions" ("user_id", "date");
-- Create index "idx_transactions_user_id" to table: "transactions"
CREATE INDEX "idx_transactions_user_id" ON "public"."transactions" ("user_id");

-- Create "transfers" table
CREATE TABLE "public"."transfers" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "from_transaction_id" uuid NOT NULL,
  "to_transaction_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_transfers_from_transaction" FOREIGN KEY ("from_transaction_id") REFERENCES "public"."transactions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_transfers_to_transaction" FOREIGN KEY ("to_transaction_id") REFERENCES "public"."transactions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_transfers_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_transfers_from_transaction" to table: "transfers"
CREATE UNIQUE INDEX "idx_transfers_from_transaction" ON "public"."transfers" ("from_transaction_id");
-- Create index "idx_transfers_to_transaction" to table: "transfers"
CREATE UNIQUE INDEX "idx_transfers_to_transaction" ON "public"."transfers" ("to_transaction_id");
-- Create index "idx_transfers_user_id" to table: "transfers"
CREATE INDEX "idx_transfers_user_id" ON "public"."transfers" ("user_id");

-- Create "budgets" table
CREATE TABLE "public"."budgets" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "category_id" uuid NOT NULL,
  "amount" numeric(15,2) NOT NULL,
  "period_type" character varying(20) NOT NULL,
  "start_date" date NOT NULL,
  "end_date" date NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_budgets_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
  CONSTRAINT "fk_budgets_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_budgets_category_id" to table: "budgets"
CREATE INDEX "idx_budgets_category_id" ON "public"."budgets" ("category_id");
-- Create index "idx_budgets_user_active" to table: "budgets"
CREATE INDEX "idx_budgets_user_active" ON "public"."budgets" ("user_id", "is_active");
-- Create index "idx_budgets_user_id" to table: "budgets"
CREATE INDEX "idx_budgets_user_id" ON "public"."budgets" ("user_id");

-- Create "budget_suggestions" table
CREATE TABLE "public"."budget_suggestions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "category_id" uuid NOT NULL,
  "suggested_amount" numeric(15,2) NOT NULL,
  "current_average" numeric(15,2) NOT NULL,
  "confidence_score" numeric(3,2) NOT NULL,
  "reason" text NULL,
  "period_type" character varying(20) NOT NULL DEFAULT 'MONTHLY',
  "generated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "applied_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_budget_suggestions_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_budget_suggestions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_budget_suggestions_generated_at" to table: "budget_suggestions"
CREATE INDEX "idx_budget_suggestions_generated_at" ON "public"."budget_suggestions" ("generated_at");
-- Create index "idx_budget_suggestions_user_id" to table: "budget_suggestions"
CREATE INDEX "idx_budget_suggestions_user_id" ON "public"."budget_suggestions" ("user_id");

-- Create "asset_snapshots" table
CREATE TABLE "public"."asset_snapshots" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "date" date NOT NULL,
  "total_assets" numeric(15,2) NOT NULL,
  "total_liabilities" numeric(15,2) NOT NULL,
  "net_worth" numeric(15,2) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_asset_snapshots_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_asset_snapshots_date" to table: "asset_snapshots"
CREATE INDEX "idx_asset_snapshots_date" ON "public"."asset_snapshots" ("date");
-- Create index "idx_asset_snapshots_user_date" to table: "asset_snapshots"
CREATE UNIQUE INDEX "idx_asset_snapshots_user_date" ON "public"."asset_snapshots" ("user_id", "date");
-- Create index "idx_asset_snapshots_user_id" to table: "asset_snapshots"
CREATE INDEX "idx_asset_snapshots_user_id" ON "public"."asset_snapshots" ("user_id");

-- Create "asset_forecasts" table
CREATE TABLE "public"."asset_forecasts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "forecast_date" date NOT NULL,
  "forecasted_amount" numeric(15,2) NOT NULL,
  "confidence_level" numeric(3,2) NOT NULL,
  "forecast_method" character varying(50) NOT NULL,
  "generated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_asset_forecasts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_asset_forecasts_generated_at" to table: "asset_forecasts"
CREATE INDEX "idx_asset_forecasts_generated_at" ON "public"."asset_forecasts" ("generated_at");
-- Create index "idx_asset_forecasts_user_date" to table: "asset_forecasts"
CREATE INDEX "idx_asset_forecasts_user_date" ON "public"."asset_forecasts" ("user_id", "forecast_date");
-- Create index "idx_asset_forecasts_user_id" to table: "asset_forecasts"
CREATE INDEX "idx_asset_forecasts_user_id" ON "public"."asset_forecasts" ("user_id");

-- Create "notification_settings" table
CREATE TABLE "public"."notification_settings" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "notification_type" character varying(50) NOT NULL,
  "is_enabled" boolean NOT NULL DEFAULT true,
  "channel" character varying(20) NOT NULL,
  "threshold_value" numeric(15,2) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notification_settings_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_notification_settings_user_id" to table: "notification_settings"
CREATE INDEX "idx_notification_settings_user_id" ON "public"."notification_settings" ("user_id");
-- Create index "idx_notification_settings_user_type" to table: "notification_settings"
CREATE UNIQUE INDEX "idx_notification_settings_user_type" ON "public"."notification_settings" ("user_id", "notification_type", "channel");

-- Create updated_at triggers for all tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON "public"."users" FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON "public"."accounts" FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON "public"."categories" FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON "public"."transactions" FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER update_budgets_updated_at BEFORE UPDATE ON "public"."budgets" FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER update_budget_suggestions_updated_at BEFORE UPDATE ON "public"."budget_suggestions" FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER update_notification_settings_updated_at BEFORE UPDATE ON "public"."notification_settings" FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();