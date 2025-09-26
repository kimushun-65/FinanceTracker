-- Add email_verified column to users table
ALTER TABLE "public"."users"
ADD COLUMN "email_verified" boolean NOT NULL DEFAULT false;

-- Create index for email_verified
CREATE INDEX "idx_users_email_verified" ON "public"."users" ("email_verified");