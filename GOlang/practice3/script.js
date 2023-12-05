const Math = require('mathjs');
const express = require('express');
const cors = require('cors');
const axios = require('axios');
const app = express();
app.use(express.json());
app.use(cors());

app.listen(8080, () => {
  console.log('Сервер запущен на порту 8080');
});

function generateShortUrl() {
  const alphabet = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890';
  let str = 'http://localhost:8080/';
  for (let i = 0; i < Math.floor(Math.random() * (3)) + 3; i++) {
    var randomNumber = Math.floor(Math.random() * (61));
    str += alphabet[randomNumber];
  }
  return str;
}

async function sendPostRequest(url, data) { // отправка пост-запроса в БД
  try {
    const response = await axios.post(url, data);
    return response.data;
  } catch (error) {
    console.error('Error:', error);
    throw new Error('Пост-запрос не отправлен: ' + error);
  }
}

async function sendGetRequest(url) { // отправка гет-запроса
  try {
    const response = await axios.get(url);
    return response.data;
  } catch (error) {
    console.error('Error:', error);
    throw new Error('Гет-запрос не отправлен: ' + error);
  }
}

app.post('/SHORT', async (req, res) => { // принятие пост-запроса от клиента
  let longUrl = req.body.longUrl;
  try {
    let response = await sendPostRequest('http://localhost:8081/insert', { key: generateShortUrl(), value: longUrl });
    res.status(200).json({ shortenedUrl: response });
  } catch (error) {
    res.status(500).json({ error: 'Не удалось сократить ссылку' });
  }
});

app.get('/:shortUrl', async (req, res) => { // переход по короткой ссылке
  let shortUrl = req.params.shortUrl;
  console.log(shortUrl);
  try {
    let response = await sendGetRequest('http://localhost:8081/get?key=' + shortUrl);
    console.log(response);
    if (response !== 'Failed to retrieve the value') {
      res.status(302).redirect(response); // перенаправляем на длинный URL
    } else {
      res.status(404).json({ error: 'Ссылка не найдена' });
    }
  } catch (error) {
    res.status(500).json({ error: 'Не удалось получить ссылку' });
  }
});
