import React from 'react';
import merge from 'lodash/merge';
import {
  Button,
  Hidden,
  Step,
  StepButton,
  Stepper,
  MobileStepper,
  StepContent,
  Paper,
  Typography, withStyles
} from '@material-ui/core';
import {
  and,
  Categorization,
  categorizationHasCategory,
  Category,
  isVisible,
  optionIs,
  RankedTester,
  rankWith,
  StatePropsOfLayout,
  uiTypeIs
} from '@jsonforms/core';
import { makeStyles, useTheme } from '@material-ui/core/styles';
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft';
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight';
import { RendererComponent, withJsonFormsLayoutProps } from '@jsonforms/react';
import {
  AjvProps,
  MaterialLayoutRenderer,
  MaterialLayoutRendererProps,
  withAjvProps
} from './layout';

export const mobileCategorizationTester: RankedTester = rankWith(
  2,
  and(
    uiTypeIs('Categorization'),
    categorizationHasCategory,
    optionIs('variant', 'stepper')
  )
);

export interface CategorizationStepperState {
  activeCategory: number;
}

export interface mobileCategorizationLayoutRendererProps
  extends StatePropsOfLayout, AjvProps {
  data: any;
}


class mobileCategorizationLayoutRenderer extends RendererComponent<mobileCategorizationLayoutRendererProps,
  CategorizationStepperState> {

  state = {
    activeCategory: 0
  };

  handleStep = (step: number) => {
    this.setState({ activeCategory: step });
  };

  render() {
    const {
      data,
      path,
      renderers,
      schema,
      uischema,
      visible,
      cells,
      config,
      ajv
    } = this.props;
    const categorization = uischema as Categorization;
    const activeCategory = this.state.activeCategory;
    const appliedUiSchemaOptions = merge({}, config, uischema.options);
    const rootStyle =  {
      width: '100%'
    }
    const buttonStyle = {
      marginTop: '8px',
      marginRight: '8px'
    }
    const actionsContainerStyle = {
      marginBottom: '16px'
    }
    const buttonWrapperStyle = {
      textAlign: 'right' as 'right',
      width: '100%',
      margin: '1em auto'
    };
    const buttonNextStyle = {
      float: 'right' as 'right'
    };
    const childProps: MaterialLayoutRendererProps = {
      elements: categorization.elements[activeCategory].elements,
      schema,
      path,
      direction: 'column',
      visible,
      renderers,
      cells
    };
    const categories = categorization.elements.filter((category: Category) =>
      isVisible(category, data, undefined, ajv)
    );
    return (
      <Hidden xsUp={!visible}>

        <div style={rootStyle}>
          <Hidden mdUp>
          <Stepper activeStep={activeCategory} orientation="vertical">
            {categories.map((e: Category, idx: number) => (
              <Step key={e.label}>
                <StepButton onClick={() => this.handleStep(idx)}>
                  {e.label}
                </StepButton>
                <StepContent>
                  <div>
                    <MaterialLayoutRenderer {...childProps} />
                  </div>
                  {!!appliedUiSchemaOptions.showNavButtons ? (
                    <div style={actionsContainerStyle}>
                      <div>
                        <Button
                          disabled={activeCategory <= 0}
                          onClick={() => this.handleStep(activeCategory - 1)}
                          style={buttonStyle}
                        >
                          Back
                        </Button>
                        <Button
                          variant="contained"
                          color="primary"
                          disabled={activeCategory >= categories.length - 1}
                          onClick={() => this.handleStep(activeCategory + 1)}
                          style={buttonStyle}
                        >
                          {/*{activeCategory === categories.length - 1 ? 'Finish' : 'Next'}*/}
                          Next
                        </Button>
                      </div>
                    </div>
                  ) : (<></>)}
                </StepContent>
              </Step>
            ))}
            {/*{activeCategory === categories.length && (*/}
            {/*  <Paper square elevation={0} className={classes.resetContainer}>*/}
            {/*    <Typography>All steps completed - you&apos;re finished</Typography>*/}
            {/*    <Button onClick={() => {}} className={buttonStyle}>*/}
            {/*      Reset*/}
            {/*    </Button>*/}
            {/*  </Paper>*/}
            {/*)}*/}
          </Stepper>
          </Hidden>
          <Hidden smDown>
          <Stepper activeStep={activeCategory} nonLinear>
          {categories.map((e: Category, idx: number) => (
            <Step key={e.label}>
              <StepButton onClick={() => this.handleStep(idx)}>
                {e.label}
              </StepButton>
            </Step>
          ))}
        </Stepper>
        <div>
          <MaterialLayoutRenderer {...childProps} />
        </div>
        { !!appliedUiSchemaOptions.showNavButtons ? (<div style={buttonWrapperStyle}>
          <Button
            style={buttonNextStyle}
            variant="contained"
            color="primary"
            disabled={activeCategory >= categories.length - 1}
            onClick={() => this.handleStep(activeCategory + 1)}
          >
            Next
          </Button>
          <Button
            style={buttonStyle}
            color="secondary"
            variant="contained"
            disabled={activeCategory <= 0}
            onClick={() => this.handleStep(activeCategory - 1)}
          >
            Previous
          </Button>
        </div>) : (<></>)}
          </Hidden>

        </div>
        {/*<MobileStepper*/}
        {/*            variant="dots"*/}
        {/*            steps={6}*/}
        {/*            position="static"*/}
        {/*            activeStep={activeCategory}*/}
        {/*            nextButton={*/}
        {/*              <Button size="small" onClick={() => this.handleStep(activeCategory + 1)} disabled={activeCategory >= categories.length - 1}>*/}
        {/*                Next*/}
        {/*              </Button>*/}
        {/*            }*/}
        {/*            backButton={*/}
        {/*              <Button size="small" onClick={() => this.handleStep(activeCategory - 1)} disabled={activeCategory <= 0}>*/}
        {/*                Back*/}
        {/*              </Button>*/}
        {/*            }*/}
        {/*/>*/}


        {/*<Stepper activeStep={activeCategory} orientation="vertical" nonLinear>*/}
        {/*  {categories.map((e: Category, idx: number) => (*/}
        {/*    <Step key={e.label}>*/}
        {/*      <StepButton onClick={() => this.handleStep(idx)}>*/}
        {/*        {e.label}*/}
        {/*      </StepButton>*/}
        {/*    </Step>*/}
        {/*  ))}*/}
        {/*</Stepper>*/}
        {/*<div>*/}
        {/*  <MaterialLayoutRenderer {...childProps} />*/}
        {/*</div>*/}
        {/*{!!appliedUiSchemaOptions.showNavButtons ? (<div style={buttonWrapperStyle}>*/}
        {/*  <Button*/}
        {/*    style={buttonNextStyle}*/}
        {/*    variant='contained'*/}
        {/*    color='primary'*/}
        {/*    disabled={activeCategory >= categories.length - 1}*/}
        {/*    onClick={() => this.handleStep(activeCategory + 1)}*/}
        {/*  >*/}
        {/*    Next*/}
        {/*  </Button>*/}
        {/*  <Button*/}
        {/*    style={buttonStyle}*/}
        {/*    color='secondary'*/}
        {/*    variant='contained'*/}
        {/*    disabled={activeCategory <= 0}*/}
        {/*    onClick={() => this.handleStep(activeCategory - 1)}*/}
        {/*  >*/}
        {/*    Previous*/}
        {/*  </Button>*/}
        {/*</div>) : (<></>)}*/}

      </Hidden>
    );
  }
}

export default withJsonFormsLayoutProps(withAjvProps(
  mobileCategorizationLayoutRenderer
));

