# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 2019_10_21_203339) do

  create_table "banks", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8", force: :cascade do |t|
    t.string "name", limit: 30, null: false
    t.string "image", default: ""
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["name"], name: "index_banks_on_name"
  end

  create_table "countries", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8", force: :cascade do |t|
    t.string "name", limit: 30, null: false
    t.string "image", default: ""
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["name"], name: "index_countries_on_name"
  end

  create_table "products", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8", force: :cascade do |t|
    t.string "title", limit: 50, null: false
    t.text "description"
    t.integer "price", default: 0, unsigned: true
    t.bigint "seller_id", null: false
    t.bigint "country_id", null: false
    t.string "image", null: false
    t.integer "status", limit: 1, default: 1, unsigned: true
    t.date "from_date", null: false
    t.date "to_date", null: false
    t.timestamp "deleted_at"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["country_id", "deleted_at"], name: "index_products_on_country_id_and_deleted_at"
    t.index ["deleted_at"], name: "index_products_on_deleted_at"
    t.index ["seller_id", "deleted_at"], name: "index_products_on_seller_id_and_deleted_at"
    t.index ["title", "country_id", "deleted_at"], name: "index_products_on_title_and_country_id_and_deleted_at"
    t.index ["title", "deleted_at"], name: "index_products_on_title_and_deleted_at"
    t.index ["title", "seller_id", "deleted_at"], name: "index_products_on_title_and_seller_id_and_deleted_at"
    t.index ["title"], name: "index_products_on_title"
  end

  create_table "user_addresses", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8", force: :cascade do |t|
    t.string "address", default: ""
    t.string "phone", limit: 20, null: false
    t.string "address_name", limit: 50, null: false
    t.bigint "user_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["user_id"], name: "index_user_addresses_on_user_id"
  end

  create_table "users", options: "ENGINE=InnoDB DEFAULT CHARSET=utf8", force: :cascade do |t|
    t.string "email"
    t.string "name", limit: 50, null: false
    t.string "phone", limit: 20, null: false
    t.string "password", null: false
    t.string "bank_name", limit: 20, null: false
    t.string "bank_account", limit: 20, null: false
    t.string "avatar", default: ""
    t.datetime "last_login_at"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["email"], name: "index_users_on_email"
    t.index ["phone"], name: "index_users_on_phone"
  end

end
