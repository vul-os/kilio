import schema from './schematest.json';
import uischema from './uischematest.json';
import {
  materialRenderers,
  materialCells,
} from '@jsonforms/material-renderers';

import './App.css';
import React, { useState } from 'react';
import { JsonForms } from '@jsonforms/react';

const initialData = {
  "provideAddress": true,
  "vegetarian": false
}

function App() {
  const [data, setData] = useState(initialData);
  return (
    <div className='App'>
      <JsonForms
        schema={schema}
        uischema={uischema}
        data={data}
        renderers={materialRenderers}
        cells={materialCells}
        onChange={({ data, _errors }) => setData(data)}
      />
    </div>
  );
}
export default App;
