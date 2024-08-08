const assetModel = require('../models/assetModel');
const xlsx = require('xlsx'); // 确保导入 xlsx 模块

const getAssets = async (req, res) => {
    const { ip, sn } = req.query;
    const field = ip ? 'ip' : 'sn_number';
    const value = ip || sn;
    try {
        const assets = await assetModel.getAssetsByField(field, value);
        res.json(assets);
    } catch (error) {
        res.status(500).json({ error: 'Database query error' });
    }
};

const createAsset = async (req, res) => {
    try {
        const assetData = req.body;
        const insertId = await assetModel.addAsset(assetData);
        res.status(201).json({ id: insertId });
    } catch (error) {
        console.error('Database insert error:', error);
        res.status(500).json({ error: 'Database insert error', details: error.message });
    }
};

const updateAsset = async (req, res) => {
    const { id } = req.params;
    try {
        await assetModel.updateAsset(id, req.body);
        res.json({ message: 'Asset updated successfully' });
    } catch (error) {
        res.status(500).json({ error: 'Database update error' });
    }
};

const deleteAsset = async (req, res) => {
    const { id } = req.params;
    try {
        await assetModel.deleteAsset(id);
        res.json({ message: 'Asset deleted successfully' });
    } catch (error) {
        res.status(500).json({ error: 'Database delete error' });
    }
};

const importAssets = async (req, res) => {
    try {
        // 打印请求的整体信息
        console.log('Request received:', req);

        // 打印文件信息，确认文件是否存在
        if (!req.file) {
            console.log('No file found in request'); // 文件不存在时的日志
            throw new Error('No file uploaded');
        }

        console.log('File received:', req.file); // 确认文件接收

        // 解析 Excel 文件
        const workbook = xlsx.read(req.file.buffer, { type: 'buffer' });
        console.log('Workbook parsed:', workbook); // 检查工作簿内容

        // 获取第一个工作表的名称
        const sheetName = workbook.SheetNames[0];
        console.log('Sheet name:', sheetName); // 检查工作表名称

        // 获取工作表并转换为 JSON 数据
        const sheet = workbook.Sheets[sheetName];
        const data = xlsx.utils.sheet_to_json(sheet);
        console.log('Data parsed from sheet:', data); // 检查解析的数据

        // 将数据插入数据库
        for (const asset of data) {
            await assetModel.addAsset(asset);
        }

        res.json({ message: 'Data imported successfully' });
    } catch (error) {
        // 捕获和打印错误信息
        console.error('Error importing data:', error.message);
        res.status(500).json({ error: 'Failed to import data' });
    }
};


module.exports = {
    getAssets,
    createAsset,
    updateAsset,
    importAssets,
    deleteAsset
};
