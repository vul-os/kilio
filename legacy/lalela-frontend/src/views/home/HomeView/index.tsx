import React, { useState } from 'react';
import type { FC } from 'react';
import { makeStyles } from '@material-ui/core';
import Page from 'src/components/Page';
import data from './data.js';
import schema from './schema.json';
import uischema from './uischema.json'
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { JsonForms } from '@jsonforms/react';
import {
  materialRenderers,
  materialCells,
} from '@jsonforms/material-renderers';
import mobileCategorizationLayoutRenderer, { mobileCategorizationTester } from './MyGroup';

const renderers = [
{ tester: mobileCategorizationTester, renderer: mobileCategorizationLayoutRenderer, },
  ...materialRenderers,
  //register custom renderers
];
const initialData = data;
const useStyles = makeStyles(() => ({
  root: {
    display:'flex', justifyContent:'center'
  },
  card: {
    marginTop: 50
  },
}));

const HomeView: FC = () => {
  const classes = useStyles();
  const [data, setData] = useState(initialData);
  return (
    <Page
      className={classes.root}
      title="Home"
    >
      <Card className={classes.card}>
        <CardContent>
          <JsonForms
            schema={schema}
            uischema={uischema}
            data={data}
            renderers={renderers}
            cells={materialCells}
            onChange={({ data }) => setData(data)}
          />
        </CardContent>

      </Card>
    </Page>
  );
};

export default HomeView;
