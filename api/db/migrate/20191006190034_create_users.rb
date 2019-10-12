class CreateUsers < ActiveRecord::Migration[5.1]
  def up
    create_table :users do |t|
      t.string :email, unique: true
      t.string :name, limit: 50, null: false
      t.string :phone, limit: 20, null: false
      t.string :password, null: false
      t.string :bank_name, limit: 20, null: false
      t.string :bank_account, limit: 20, null: false
      t.string :avatar, default: ""
      t.datetime :last_login_at

      t.timestamps

      t.index :email
      t.index :phone
    end
  end

  def down
    drop_table :users
  end
end
