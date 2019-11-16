class AddReceiptProofAndPaymentMethodToInvoices < ActiveRecord::Migration[5.1]
  def up
    change_table :invoices do |t|
      t.string :receipt_proof, limit: 255, default: ""
      t.string :payment_method, limit: 50, null: false
    end
  end

  def down
    change_table :invoices do |t|
      t.remove :receipt_proof
      t.remove :payment_method
    end
  end
end
