import React from 'react';
import {
  and,
  RankedTester,
  rankWith,
  scopeEndsWith,
} from '@jsonforms/core';
import { DropzoneArea } from 'material-ui-dropzone';
import { withJsonFormsControlProps } from '@jsonforms/react';


interface FileUploadControlProps {
    data: any;
    path: string;
  }

export const UploadTester: RankedTester = rankWith(
  3,
  and(
    scopeEndsWith('uploadDoc')
  )
);


const FileUploadControl = ({ data, path }: FileUploadControlProps) => (
    <DropzoneArea
        onChange={(files) => console.log('Files:', files)}
    />
  );

export default withJsonFormsControlProps(FileUploadControl)