class CreateTransactions < ActiveRecord::Migration[5.1]
  def up
    create_table :transactions do |t|
      t.bigint   :product_id, null: false
      t.bigint   :buyer_id, null: false
      t.bigint   :seller_id, null: false
      t.bigint   :buyer_address_id, null: false
      t.integer  :quantity, unsigned: true, limit: 2, default: 1
      t.string   :notes, limit: 200, default: ""
      t.integer  :total_price, unsigned: true, limit: 4, null: false
      t.integer  :status, unsigned: true, limit: 1, default: 0
      t.datetime :paid_at, null: true
      t.datetime :finished_at, null: true

      t.timestamps

      t.index :product_id
      t.index :buyer_id
      t.index :seller_id
      t.index :buyer_address_id
    end
  end

  def down
    drop_table :transactions
  end
end
