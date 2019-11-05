class CreateInvoices < ActiveRecord::Migration[5.1]
  def up
    create_table :invoices do |t|
      t.bigint   :transaction_id, null: false
      t.string   :invoice_code, null: false
      t.integer  :coded_price, limit: 4, null: false
      t.integer  :status, limit: 1, default: 0

      t.timestamps

      t.index :transaction_id
      t.index :invoice_code
      t.index :status
    end
  end

  def down
    drop_table :invoices
  end
end
