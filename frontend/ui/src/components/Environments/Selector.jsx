import { useEffect, useState } from "react";
import { Radio } from "antd";


const EnvironmentsSelector = ({project, environment, onEnvChange, style, defaultEnvironment="prod"}) => {
    const [environments, setEnvironments] = useState([])
    const [currentEnvironment, setCurrentEnvironment] = useState(environment)

    useEffect(() => {
        fetch(`/api/v1/projects/${project}/envs`)
            .then((response) => {
                if (!response.ok) throw new Error("err");
                return response.json();
            })
            .then((data) => {
                setEnvironments(data.data.environments)

                if (environment === undefined) {
                    setCurrentEnvironment(defaultEnvironment)
                    onEnvChange(defaultEnvironment)
                }
            });
    }, []);

    const handleCurrentEnvironmentChange = (env) => {
        setCurrentEnvironment(env)
        onEnvChange(env)
    }

    return (
        <Radio.Group
            value={currentEnvironment}
            onChange={(e) => handleCurrentEnvironmentChange(e.target.value)}
            optionType="button"
            buttonStyle="solid"
            style={style}
        >
            {environments.map(env => (
                <Radio.Button key={env.name} value={env.name}>{env.name}</Radio.Button>
            ))}
        </Radio.Group>
    )
}

export default EnvironmentsSelector;