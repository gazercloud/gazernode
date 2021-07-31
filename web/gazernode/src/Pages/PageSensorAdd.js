import React from 'react';
import Button from "@material-ui/core/Button";
import Stepper from "@material-ui/core/Stepper";
import Step from "@material-ui/core/Step";
import StepLabel from "@material-ui/core/StepLabel";
import Paper from "@material-ui/core/Paper";
import WidgetSensorTypeSelect from "../Widgets/WidgetSensorTypeSelect";
import WidgetNewSensorName from "../Widgets/WidgetNewSensorName";
import WidgetSensorsConfiguration from "../Widgets/WidgetSensorConfiguration";

export default function PageSensorAdd(props) {

    const getStepContent = (step) => {
        switch (step) {
            case 0:
                return 'Select sensor type (' + sensorType + ')';
            case 1:
                return 'Enter sensor name (' + sensorName + ')';
            case 2:
                return 'Configure sensor';
            default:
                return 'Unknown step';
        }
    }

    const [activeStep, setActiveStep] = React.useState(0);

    const [sensorType, setSensorType] = React.useState("")
    const [sensorName, setSensorName] = React.useState("")
    const [sensorConfig, setSensorConfig] = React.useState("")

    const [messageError, setMessageError] = React.useState("")

    const handleNext = () => {
        if (activeStep < 2)
            setActiveStep(activeStep + 1);
    };

    const handleBack = () => {
        if (activeStep > 0)
            setActiveStep(activeStep - 1);
    };

    const stepSelectSensorType = () => {
        if (activeStep !== 0) {
            return (<div/>);
        }

        return (
            <WidgetSensorTypeSelect onSelected={(sType) => {
                setSensorType(sType)
                handleNext()
            }
            } />
        )
    }

    const stepEnterSensorName = () => {
        if (activeStep !== 1) {
            return (<div/>);
        }
        return (
            <WidgetNewSensorName Cancel={handleBack} OK={(sName) => {
                setSensorName(sName)
                handleNext()
            }} />
        )
    }

    const stepSensorConfiguration = () => {
        if (activeStep !== 2) {
            return (<div/>);
        }
        return (
            <WidgetSensorsConfiguration />
        )
    }

    const addSensor = (type, name) => {
        let req = {
            "fn": "host_add_sensor",
            "type": type,
            "name": name
        }
        fetch("/api/request?request=" + JSON.stringify(req))
            .then(res => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
                //setProcessingLoadParameter(false)
            })
            .catch((err) => {
                //setProcessingLoadParameter(false)
                //setErrorMessage("Unknown error")
            })
    }


    return (
        <div>
            <Stepper activeStep={activeStep}>
                <Step key="sensor_add_step_0">
                    <StepLabel>{getStepContent(0)}</StepLabel>
                </Step>
                <Step key="sensor_add_step_1">
                    <StepLabel>{getStepContent(1)}</StepLabel>
                </Step>
                <Step key="sensor_add_step_2">
                    <StepLabel>{getStepContent(2)}</StepLabel>
                </Step>
            </Stepper>

            <Paper>
                <div style={{margin: "10px"}}>
                    {stepSelectSensorType()}
                    {stepEnterSensorName()}
                    {stepSensorConfiguration()}
                </div>
            </Paper>


            <Button
                onClick={() => {
                    addSensor(sensorType, sensorName)
                    props.onComplete()
                }}>
                OK
            </Button>
        </div>
    );
}
