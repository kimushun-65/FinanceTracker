// FinanceTracker Database Schema

table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "auth0_id" {
    null = false
    type = varchar(255)
  }
  column "email" {
    null = false
    type = varchar(255)
  }
  column "name" {
    null = false
    type = varchar(255)
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_auth0_id" {
    unique = true
    columns = [column.auth0_id]
  }
  index "idx_users_email" {
    unique = true
    columns = [column.email]
  }
}

table "accounts" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "name" {
    null = false
    type = varchar(255)
  }
  column "type" {
    null = false
    type = varchar(50)
    comment = "CASH, BANK, CREDIT_CARD, INVESTMENT, LOAN, OTHER"
  }
  column "balance" {
    null = false
    type = decimal(15, 2)
    default = 0.00
  }
  column "currency" {
    null = false
    type = varchar(3)
    default = "JPY"
  }
  column "is_active" {
    null = false
    type = boolean
    default = true
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_accounts_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  index "idx_accounts_user_id" {
    columns = [column.user_id]
  }
}

table "category_masters" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "name" {
    null = false
    type = varchar(100)
  }
  column "type" {
    null = false
    type = varchar(20)
    comment = "INCOME, EXPENSE"
  }
  column "icon" {
    null = true
    type = varchar(50)
  }
  column "color" {
    null = true
    type = varchar(7)
    comment = "Hex color code"
  }
  column "display_order" {
    null = false
    type = integer
    default = 0
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_category_masters_type" {
    columns = [column.type]
  }
}

table "categories" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "category_master_id" {
    null = false
    type = uuid
  }
  column "custom_name" {
    null = true
    type = varchar(100)
  }
  column "is_active" {
    null = false
    type = boolean
    default = true
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_categories_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  foreign_key "fk_categories_master" {
    columns = [column.category_master_id]
    ref_columns = [table.category_masters.column.id]
    on_delete = RESTRICT
  }
  index "idx_categories_user_id" {
    columns = [column.user_id]
  }
  index "idx_categories_master_id" {
    columns = [column.category_master_id]
  }
  index "idx_categories_user_master" {
    unique = true
    columns = [column.user_id, column.category_master_id]
  }
}

table "transactions" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "account_id" {
    null = false
    type = uuid
  }
  column "category_id" {
    null = false
    type = uuid
  }
  column "amount" {
    null = false
    type = decimal(15, 2)
  }
  column "type" {
    null = false
    type = varchar(20)
    comment = "INCOME, EXPENSE, TRANSFER"
  }
  column "date" {
    null = false
    type = date
  }
  column "description" {
    null = true
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_transactions_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  foreign_key "fk_transactions_account" {
    columns = [column.account_id]
    ref_columns = [table.accounts.column.id]
    on_delete = RESTRICT
  }
  foreign_key "fk_transactions_category" {
    columns = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_delete = RESTRICT
  }
  index "idx_transactions_user_id" {
    columns = [column.user_id]
  }
  index "idx_transactions_account_id" {
    columns = [column.account_id]
  }
  index "idx_transactions_category_id" {
    columns = [column.category_id]
  }
  index "idx_transactions_date" {
    columns = [column.date]
  }
  index "idx_transactions_user_date" {
    columns = [column.user_id, column.date]
  }
}

table "transfers" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "from_transaction_id" {
    null = false
    type = uuid
  }
  column "to_transaction_id" {
    null = false
    type = uuid
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_transfers_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  foreign_key "fk_transfers_from_transaction" {
    columns = [column.from_transaction_id]
    ref_columns = [table.transactions.column.id]
    on_delete = CASCADE
  }
  foreign_key "fk_transfers_to_transaction" {
    columns = [column.to_transaction_id]
    ref_columns = [table.transactions.column.id]
    on_delete = CASCADE
  }
  index "idx_transfers_user_id" {
    columns = [column.user_id]
  }
  index "idx_transfers_from_transaction" {
    unique = true
    columns = [column.from_transaction_id]
  }
  index "idx_transfers_to_transaction" {
    unique = true
    columns = [column.to_transaction_id]
  }
}

table "budgets" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "category_id" {
    null = false
    type = uuid
  }
  column "amount" {
    null = false
    type = decimal(15, 2)
  }
  column "period_type" {
    null = false
    type = varchar(20)
    comment = "MONTHLY, YEARLY"
  }
  column "start_date" {
    null = false
    type = date
  }
  column "end_date" {
    null = true
    type = date
  }
  column "is_active" {
    null = false
    type = boolean
    default = true
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_budgets_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  foreign_key "fk_budgets_category" {
    columns = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_delete = RESTRICT
  }
  index "idx_budgets_user_id" {
    columns = [column.user_id]
  }
  index "idx_budgets_category_id" {
    columns = [column.category_id]
  }
  index "idx_budgets_user_active" {
    columns = [column.user_id, column.is_active]
  }
}

table "budget_suggestions" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "category_id" {
    null = false
    type = uuid
  }
  column "suggested_amount" {
    null = false
    type = decimal(15, 2)
  }
  column "current_average" {
    null = false
    type = decimal(15, 2)
  }
  column "confidence_score" {
    null = false
    type = decimal(3, 2)
    comment = "0.00 to 1.00"
  }
  column "reason" {
    null = true
    type = text
  }
  column "period_type" {
    null = false
    type = varchar(20)
    default = "MONTHLY"
  }
  column "generated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "applied_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_budget_suggestions_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  foreign_key "fk_budget_suggestions_category" {
    columns = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_delete = CASCADE
  }
  index "idx_budget_suggestions_user_id" {
    columns = [column.user_id]
  }
  index "idx_budget_suggestions_generated_at" {
    columns = [column.generated_at]
  }
}

table "asset_snapshots" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "date" {
    null = false
    type = date
  }
  column "total_assets" {
    null = false
    type = decimal(15, 2)
  }
  column "total_liabilities" {
    null = false
    type = decimal(15, 2)
  }
  column "net_worth" {
    null = false
    type = decimal(15, 2)
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_asset_snapshots_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  index "idx_asset_snapshots_user_id" {
    columns = [column.user_id]
  }
  index "idx_asset_snapshots_date" {
    columns = [column.date]
  }
  index "idx_asset_snapshots_user_date" {
    unique = true
    columns = [column.user_id, column.date]
  }
}

table "asset_forecasts" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "forecast_date" {
    null = false
    type = date
  }
  column "forecasted_amount" {
    null = false
    type = decimal(15, 2)
  }
  column "confidence_level" {
    null = false
    type = decimal(3, 2)
    comment = "0.00 to 1.00"
  }
  column "forecast_method" {
    null = false
    type = varchar(50)
    comment = "LINEAR, EXPONENTIAL, SEASONAL"
  }
  column "generated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_asset_forecasts_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  index "idx_asset_forecasts_user_id" {
    columns = [column.user_id]
  }
  index "idx_asset_forecasts_generated_at" {
    columns = [column.generated_at]
  }
  index "idx_asset_forecasts_user_date" {
    columns = [column.user_id, column.forecast_date]
  }
}

table "notification_settings" {
  schema = schema.public
  column "id" {
    null = false
    type = uuid
    default = sql("gen_random_uuid()")
  }
  column "user_id" {
    null = false
    type = uuid
  }
  column "notification_type" {
    null = false
    type = varchar(50)
    comment = "BUDGET_EXCEEDED, PAYMENT_DUE, LOW_BALANCE, MONTHLY_SUMMARY"
  }
  column "is_enabled" {
    null = false
    type = boolean
    default = true
  }
  column "channel" {
    null = false
    type = varchar(20)
    comment = "EMAIL, PUSH, SMS"
  }
  column "threshold_value" {
    null = true
    type = decimal(15, 2)
  }
  column "created_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null = false
    type = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_notification_settings_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  index "idx_notification_settings_user_id" {
    columns = [column.user_id]
  }
  index "idx_notification_settings_user_type" {
    unique = true
    columns = [column.user_id, column.notification_type, column.channel]
  }
}

// Schema for public schema
schema "public" {
}