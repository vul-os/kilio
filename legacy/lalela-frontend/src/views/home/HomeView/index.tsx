import React, { useState } from 'react';
import type { FC } from 'react';
import { makeStyles } from '@material-ui/core';
import Page from 'src/components/Page';
import Hero from './Hero';
import Features from './Features';
import Testimonials from './Testimonials';
import CTA from './CTA';
import FAQS from './FAQS';

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';

import { JsonForms } from '@jsonforms/react';
import {
  materialRenderers,
  materialCells,
} from '@jsonforms/material-renderers';
import { person } from '@jsonforms/examples';
const initialData = person.data;
const useStyles = makeStyles(() => ({
  root: {
    display:'flex', justifyContent:'center'
  },
  card: {
    marginTop: 50
  },
}));

const scheme = {
  "type": "object",
  "properties": {
    "WitnessOrVictim": {
      "type": "string",
      "enum": [
        "Victim",
        "Witness"
      ]
    },
    "victimAwareness": {
      "type": "string",
      "enum": [
        "Yes",
        "No"
      ]
    },
    "WitnessFlow": {
      "type": "object",
      "properties": {
        "WitnessWhen": {
          "type": "object",
          "properties": {
            "DateRecall": {
              "type": "string",
              "enum": [
                "Yes",
                "No"
              ]
            },
            "IncidentDate":{
              "type": "string",
              "format": "date"
            }
          }
        }
      }
    }
  }
}

const uischeme = {
  "type": "Categorization",
  "elements": [
    {
      "type": "Category",
      "label": "Basic Intro",
      "elements": [
        {
          "type": "VerticalLayout",
          "elements": [
            {
              "type": "Control",
              "label": "Please indicate if you are reporting on behalf of someone",
              "scope": "#/properties/WitnessOrVictim",
              "options": {
                "format": "radio"
              }

            },
            {"type": "Control",
              "label": "",
              "scope": "#/properties/WitnessFlow/properties/WitnessWhen/properties/IncidentDate",
              "rule": {
                "effect": "SHOW",
                "condition": {
                  "scope": "#/properties/WitnessFlow/properties/WitnessWhen/properties/DateRecall",
                  "schema": {
                    "const": "Yes"
                  }

                }
              }}
          ]
        }

      ]},
    {
      "type": "Category",
      "label": "Parties involved"},
    {
      "type": "Category",
      "label": "Extent of incident"},
    {
      "type": "Category",
      "label": "Emotional Evaluation"},
    {
      "type": "Category",
      "label": "Supporting documents"}
  ],
  "options": {
    "variant": "stepper",
    "showNavButtons": true
  }
}

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
            schema={scheme}
            uischema={uischeme}
            data={data}
            renderers={materialRenderers}
            cells={materialCells}
            onChange={({ data }) => setData(data)}
          />
        </CardContent>

      </Card>
    </Page>
  );
};

export default HomeView;
