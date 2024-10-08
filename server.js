const express = require('express');
const cors = require('cors');
const assetRoutes = require('./routes/assetRoutes');
const upload = require('./middleware/upload');
const xlsx = require('xlsx');
const assetModel = require('./models/assetModel');
const initDb = require('./initDb');
require('dotenv').config();

const app = express();
app.use(cors());
app.use(express.json());

app.use('/api', assetRoutes);

// 初始化数据库表
initDb();

// 导入Excel数据
app.post('/api/assets/import', upload.single('file'), async (req, res) => {
    try {
        const workbook = xlsx.read(req.file.buffer, { type: 'buffer' });
        const sheetName = workbook.SheetNames[0];
        const sheet = workbook.Sheets[sheetName];
        const data = xlsx.utils.sheet_to_json(sheet);

        for (const asset of data) {
            await assetModel.addAsset(asset);
        }

        res.json({ message: 'Data imported successfully' });
    } catch (error) {
        res.status(500).json({ error: 'Failed to import data' });
    }
});

// 导出Excel数据
app.get('/api/assets/export', async (req, res) => {
    try {
        const assets = await assetModel.getAssetsByField('1', '1');
        const worksheet = xlsx.utils.json_to_sheet(assets);
        const workbook = xlsx.utils.book_new();
        xlsx.utils.book_append_sheet(workbook, worksheet, 'Assets');

        res.setHeader(
            'Content-Disposition',
            'attachment; filename="assets.xlsx"'
        );
        res.send(xlsx.write(workbook, { type: 'buffer', bookType: 'xlsx' }));
    } catch (error) {
        res.status(500).json({ error: 'Failed to export data' });
    }
});
const corsOptions = {
    origin: 'http://193.26.156.158:3000', // 替换为你的前端应用的地址
    optionsSuccessStatus: 200 // 一些旧版浏览器对204的响应处理不佳
};

app.use(cors(corsOptions));


const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
});
