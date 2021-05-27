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
        },
        "WitnessWhere": {
          "type": "object",
          "properties": {
            "Where": {
              "type": "string"
            }
          }
        }
      }
    }
  },
  "required": [
    "WitnessOrVictim"
  ]
}

const uischeme = {
  "type": "Categorization",
  "elements": [
    {
      "type": "Category",
      "label": "Basic Intro",
      "elements": [
        {
          "type": "Group",
          "label": "Before we begin please let us know if you are reporting an incident agaisnt yourself, or on behalf of another person",
          "elements": [
            {
              "type": "VerticalLayout",
              "elements": [
                {
                  "type": "Control",
                  "scope": "#/properties/WitnessOrVictim",
                  "label": "Please indicate if you are reporting on behalf of someone",
                  "options": {
                    "format": "radio"
                  }
                },
                {
                  "type": "Control",
                  "label": "Does the victim know you are filing this report?",
                  "scope": "#/properties/victimAwareness",
                  "options": {
                    "format": "radio"
                  },
                  "rule": {
                    "effect": "SHOW",
                    "condition": {
                      "scope": "#/properties/WitnessOrVictim",
                      "schema": {
                        "const": "Witness"
                      }

                    }
                  }
                }
              ]
            }
          ]
        },

      ]},
    {
      "type": "Category",
      "label": "Time, Date & Location",
      "elements": [
        {
          "type": "VerticalLayout",
          "elements": [
            {"type": "Label",
              "text": "Thank you for taking the time to to report what you have witnessed. This section will focus on the 'where' and 'when' of the incident LOOOONGER SENNNNTENCE"},

            {
              "type": "HorizontalLayout",
              "elements": [
                {"type": "Control",
                  "label": "Can you recall the Date of the incident?",
                  "scope": "#/properties/WitnessFlow/properties/WitnessWhen/properties/DateRecall",
                  "options": {
                    "format": "radio"
                  }
                },
                {"type": "Control",
                  "label": "Please indicate the date",
                  "scope": "#/properties/WitnessFlow/properties/WitnessWhen/properties/IncidentDate",
                  "rule": {
                    "effect": "SHOW",
                    "condition": {
                      "scope": "#/properties/WitnessFlow/properties/WitnessWhen/properties/DateRecall",
                      "schema": {
                        "const": "Yes"
                      }

                    }
                  }}]},



            {"type":"Control",
              "label": "Where did this incident take place? Any and all information could be imprtant, please give as much detail as possible",
              "scope": "#/properties/WitnessFlow/properties/WitnessWhere/properties/Where" }
          ]
        }

      ],
      "rule":{
        "effect":"SHOW",
        "condition":{
          "scope": "#/properties/WitnessFlow/properties/WitnessWhen/properties/DateRecall",
          "schema": {
            "enum": [
              "Yes",
              "No"
            ]
          }

        }
      }
    },
    {
      "type": "Category",
      "label": "Parties involved",
      "rule":{
        "effect":"SHOW",
        "condition":{
          "scope": "#/properties/WitnessOrVictim",
          "schema": {
            "enum": [
              "Victim",
              "Witness"
            ]
          }

        }
      }}

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
