class CreateTransactionShippings < ActiveRecord::Migration[5.1]
  def up
    create_table :transaction_shippings do |t|
      t.bigint :transaction_id, null: false
      t.string :awb_number, limit: 100, default: ""
      t.string :courier, limit: 64, default: ""

      t.timestamps

      t.index :transaction_id
      t.index :courier
    end
  end

  def down
    drop_table :transaction_shippings
  end
end
