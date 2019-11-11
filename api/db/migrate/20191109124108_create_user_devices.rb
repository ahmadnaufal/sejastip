class CreateUserDevices < ActiveRecord::Migration[5.1]
  def up
    create_table :user_devices do |t|
      t.string :device_id, limit: 270, null: false
      t.string :platform, limit: 50, null: false
      t.string :user_agent, limit: 255, default: ""
      t.bigint :user_id, unsigned: false, null: false

      t.timestamps

      t.index [:user_id, :platform]
      t.index :platform
      t.index :user_id
    end
  end

  def down
    drop_table :user_devices
  end
end
