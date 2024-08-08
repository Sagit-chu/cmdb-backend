const db = require('./config/db');

const initDb = async () => {
    const createTableQuery = `
    CREATE TABLE IF NOT EXISTS cmdb_assets (
      id INT AUTO_INCREMENT PRIMARY KEY,
      ip VARCHAR(15) NOT NULL,
      application_system VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
      application_manager VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
      overall_manager VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
      is_virtual_machine BOOLEAN,
      resource_pool VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
      data_center VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
      rack_location VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
      sn_number VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
      out_of_band_ip VARCHAR(15),
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
  `;

    try {
        await db.query(createTableQuery);
        console.log('Database tables initialized successfully.');
    } catch (error) {
        console.error('Error initializing database tables:', error);
    }
};

module.exports = initDb;
