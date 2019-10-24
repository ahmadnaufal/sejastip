class CreateProducts < ActiveRecord::Migration[5.1]
  def up
    create_table :products do |t|
      t.string  :title, limit: 50, null: false
      t.text    :description, limit: 300
      t.integer :price, unsigned: true, limit: 4, default: 0
      t.bigint  :seller_id, null: false
      t.bigint  :country_id, null: false
      t.string  :image, null: false
      t.integer :status, unsigned: true, limit: 1, default: 1
      t.date    :from_date, null: false
      t.date    :to_date, null: false
      t.timestamp :deleted_at

      t.timestamps

      t.index [:title, :country_id, :deleted_at]
      t.index [:title, :seller_id, :deleted_at]
      t.index [:title, :deleted_at]
      t.index [:seller_id, :deleted_at]
      t.index [:country_id, :deleted_at]
      t.index :title
      t.index :deleted_at
    end
  end

  def down
    drop_table :products
  end
end
