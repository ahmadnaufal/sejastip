class AddNewColumnsToInvoicesAndTransactions < ActiveRecord::Migration[5.1]
  def up
    change_table :invoices do |t|
      t.string :receipt_proof, limit: 255, default: ""
      t.string :payment_method, limit: 50, null: false
      t.datetime :paid_at, null: true
    end

    change_table :transactions do |t|
      t.bigint :invoice_id, null: true
    end
  end

  def down
    change_table :invoices do |t|
      t.remove :receipt_proof
      t.remove :payment_method
      t.remove :paid_at
    end

    change_table :transactions do |t|
      t.remove :invoice_id
    end
  end
end
