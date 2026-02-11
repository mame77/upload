import { useState, ChangeEvent } from 'react';

const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8000';

const App = () => {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState<string>('');

  const hoge = (e: ChangeEvent<HTMLInputElement>) => {
    setFile(e.target.files?.[0] ?? null);
  };

  const moge = async () => {
    if (!file) return;

    const contentType = file.type || 'application/octet-stream';
    setStatus('presign...');
    const presignRes = await fetch(`${API_BASE}/presign`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        fileName: file.name,
        contentType,
      }),
    });
    if (!presignRes.ok) {
      setStatus('presign failed');
      return;
    }
    const presign = await presignRes.json();

    setStatus('uploading...');
    const putRes = await fetch(presign.url, {
      method: presign.method ?? 'PUT',
      headers: presign.headers ?? { 'Content-Type': contentType },
      body: file,
    });
    if (!putRes.ok) {
      setStatus('upload failed');
      return;
    }

    setStatus('saving...');
    const completeRes = await fetch(`${API_BASE}/upload/complete`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        objectKey: presign.objectKey,
        fileName: file.name,
        contentType,
        size: file.size,
      }),
    });
    if (!completeRes.ok) {
      setStatus('save failed');
      return;
    }
    setStatus('done');
  };

  return (
    <>
      <input type="file" onChange={hoge} />
      <button onClick={moge}>send</button>
      <div>{status}</div>
    </>
  );
};

export default App;
