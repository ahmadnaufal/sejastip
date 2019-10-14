class CreateBanks < ActiveRecord::Migration[5.1]
  def up
    create_table :banks do |t|
      t.string :name, limit: 30, null: false
      t.string :image, default: ""

      t.timestamps

      t.index :name
    end
  end

  def down
    drop_table :banks
  end
end
