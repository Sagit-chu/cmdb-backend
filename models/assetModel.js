const db = require('../config/db');

const getAssetsByField = async (field, value) => {
    const [rows] = await db.query(`SELECT * FROM cmdb_assets WHERE ${field} = ?`, [value]);
    return rows;
};

const addAsset = async (assetData) => {
    const [result] = await db.query('INSERT INTO cmdb_assets SET ?', assetData);
    return result.insertId;
};

const updateAsset = async (id, assetData) => {
    try {
        await db.query('UPDATE cmdb_assets SET ? WHERE id = ?', [assetData, id]);
    } catch (error) {
        console.error('SQL Update Error:', error); // 输出错误信息
        throw error;
    }
};


const deleteAsset = async (id) => {
    await db.query('DELETE FROM cmdb_assets WHERE id = ?', [id]);
};

module.exports = {
    getAssetsByField,
    addAsset,
    updateAsset,
    deleteAsset
};
