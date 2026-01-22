import { useState, ChangeEvent } from 'react';

const App = () => {
  const [file, setFile] = useState<File | null>(null);

  const hoge = (e: ChangeEvent<HTMLInputElement>) => {
    setFile(e.target.files?.[0] ?? null);
  };

  const moge = async () => {
    if (!file) return;

    const formData = new FormData();
    formData.append('image', file);

    await fetch('http://localhost:8000/upload', {
      method: 'POST',
      body: formData,
    });
  };

  return (
    <>
      <input type="file" onChange={hoge} />
      <button onClick={moge}>send</button>
    </>
  );
};

export default App;
