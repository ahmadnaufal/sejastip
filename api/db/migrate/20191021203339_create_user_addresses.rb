class CreateUserAddresses < ActiveRecord::Migration[5.1]
  def up
    create_table :user_addresses do |t|
      t.string :address, default: ""
      t.string :phone, limit: 20, null: false
      t.string :address_name, limit: 50, null: false
      t.bigint :user_id, null: false

      t.timestamps

      t.index :user_id
    end
  end

  def down
    drop_table :user_addresses
  end
end
