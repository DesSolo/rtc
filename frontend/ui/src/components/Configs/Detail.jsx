import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {useOutletContext, useParams} from "react-router-dom";
import {
    Card,
    Button,
    Input,
    InputNumber,
    Switch,
    Typography,
    Row,
    Col,
    Divider,
    notification,
} from "antd";
import { LockOutlined, SearchOutlined, ReloadOutlined } from "@ant-design/icons";
import EnvironmentsSelector from "../Environments/Selector";
import {fetchWithAuth} from "../../utils/fetchWithAuth.js";

const { Text } = Typography;

// эскейп для регулярного выражения
const escapeRegExp = (s = "") => s.replace(/[.*+?^${}()|[\\]\\\\]/g, "\\$&");

const Highlight = ({ text = "", query = "" }) => {
    if (!query) return <>{text}</>;
    const regex = new RegExp(`(${escapeRegExp(query)})`, "gi");
    const parts = text.split(regex);
    return (
        <>
            {parts.map((part, i) =>
                    part.toLowerCase() === query.toLowerCase() ? (
                        <span key={i} style={{ backgroundColor: "yellow" }}>
            {part}
          </span>
                    ) : (
                        part
                    )
            )}
        </>
    );
};

const ConfigRow = ({ title, description, control, locked = false, query, isChanged }) => (
    <Row align="middle" justify="flex-start" style={{ padding: "12px 0" }}>
        <Col flex="auto" style={{ maxWidth: 500 }}>
            <Text strong>
                <Highlight text={title} query={query} />
            </Text>
            <br />
            <Text type="secondary">{description}</Text>
        </Col>
        <Col style={{ display: "flex", alignItems: "center", gap: 8 }}>
            <div style={isChanged ? { border: "2px solid #1890ff", borderRadius: "4px" } : {}}>
                {control}
            </div>
            {locked && <LockOutlined style={{ color: "#999" }} />}
        </Col>
    </Row>
);

const ConfigsList = () => {
    const { setTitle } = useOutletContext();
    const { project, environment, release } = useParams();
    const groupRefs = useRef({});
    const highlightTimer = useRef(null);
    const [api, contextHolder] = notification.useNotification();

    const [configs, setConfigs] = useState([]);
    const [loading, setLoading] = useState(false);
    const [currentEnv, setCurrentEnv] = useState(environment);
    const [originalValues, setOriginalValues] = useState(new Map());
    const [modifiedValues, setModifiedValues] = useState(new Map());
    const [highlighted, setHighlighted] = useState(null);
    const [filter, setFilter] = useState("");

    // fetchConfigs — мемоизированная, чтобы не пересоздавать в эффектах
    const fetchConfigs = useCallback(async (env) => {
        try {
            const resp = await fetchWithAuth(`/api/v1/projects/${project}/envs/${env}/releases/${release}/configs`);
            if (!resp.ok) throw new Error(`status ${resp.status}`);
            const data = await resp.json();
            const configsData = data.data?.configs || [];
            setConfigs(configsData);

            const newOriginal = new Map();
            configsData.forEach((cfg) => newOriginal.set(cfg.key, cfg.value));
            setOriginalValues(newOriginal);
            setModifiedValues(new Map());
            setCurrentEnv(env);
        } catch (err) {
            api.error({ message: "Failed to load configs", description: String(err) });
        }
    }, [project, release, api]);

    // Изначальная загрузка и перезагрузка при смене проекта/release
    useEffect(() => {
        fetchConfigs(currentEnv);
        setTitle(`Configs: ${project}/${release}`)
    }, [fetchConfigs]);

    const handleChange = useCallback((key, value) => {
        // приводим к строке для хранения
        const newValue = value === undefined || value === null ? "" : String(value);
        const originalValue = originalValues.get(key) ?? "";

        setConfigs((prev) => prev.map((cfg) => (cfg.key === key ? { ...cfg, value: newValue } : cfg)));

        setModifiedValues((prev) => {
            const next = new Map(prev);
            if (newValue === originalValue) {
                next.delete(key);
            } else {
                next.set(key, { oldValue: originalValue, newValue });
            }
            return next;
        });
    }, [originalValues]);

    const handleSave = useCallback(async () => {
        const changes = {};
        for (const [key, val] of modifiedValues.entries()) changes[key] = val.newValue;

        if (Object.keys(changes).length === 0) return;

        try {
            const resp = await fetchWithAuth(`/api/v1/projects/${project}/envs/${currentEnv}/releases/${release}/configs`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(changes),
            });
            if (!resp.ok) {
                api.error({ message: "Failed to update", description: String(resp.status) });
                return;
            }
            api.success({ message: "Updated successfully" });
            await fetchConfigs(currentEnv);
        } catch (err) {
            api.error({ message: "Save error", description: String(err) });
        }
    }, [modifiedValues, project, currentEnv, release, api, fetchConfigs]);

    const grouped = useMemo(() => {
        return configs.reduce((acc, cfg) => {
            const g = (cfg.group || "UNSPECIFIED").toUpperCase();
            if (!acc[g]) acc[g] = [];
            acc[g].push(cfg);
            return acc;
        }, {});
    }, [configs]);

    const filteredGroups = useMemo(() => {
        if (!filter) return grouped;
        const lq = filter.toLowerCase();
        return Object.entries(grouped).reduce((acc, [group, items]) => {
            const matched = items.filter((cfg) =>
                cfg.key.toLowerCase().includes(lq) || (cfg.usage || "").toLowerCase().includes(lq)
            );
            if (matched.length) acc[group] = matched;
            return acc;
        }, {});
    }, [grouped, filter]);

    const renderControl = (cfg) => {
        const commonProps = { disabled: !cfg.writable, style: { minWidth: 120 } };

        switch (cfg.value_type) {
            case "string":
                return (
                    <Input
                        {...commonProps}
                        value={cfg.value ?? ""}
                        onChange={(e) => handleChange(cfg.key, e.target.value)}
                        allowClear
                    />
                );
            case "bool":
                return (
                    <Switch
                        {...commonProps}
                        checked={String(cfg.value) === "true"}
                        onChange={(val) => handleChange(cfg.key, val ? "true" : "false")}
                    />
                );
            case "int":
            case "int64":
            case "uint":
            case "uint64":
            case "float":
            case "float64":
                return (
                    <InputNumber
                        {...commonProps}
                        value={cfg.value === "" ? undefined : Number(cfg.value)}
                        onChange={(val) => handleChange(cfg.key, val)}
                        parser={(v) => v}
                    />
                );
            default:
                return <Text type="secondary">Unsupported type: {cfg.value_type}</Text>;
        }
    };

    const handleScrollToGroup = (group) => {
        const el = groupRefs.current[group];
        if (el && el.scrollIntoView) {
            el.scrollIntoView({ behavior: "smooth", block: "start" });
            setHighlighted(group);
            if (highlightTimer.current) clearTimeout(highlightTimer.current);
            highlightTimer.current = setTimeout(() => setHighlighted(null), 1400);
        }
    };

    const handleReload = async () => {
        setLoading(true)
        await fetchConfigs(currentEnv);
        api.success({ message: "Reloaded" });
        setLoading(false)
    };

    return (
        <>
            {contextHolder}
            <div style={{ display: "flex", gap: 12, marginBottom: 12 }}>
                <Button type="dashed" icon={<ReloadOutlined />} onClick={handleReload} loading={loading} />
                <EnvironmentsSelector
                    project={project}
                    environment={currentEnv}
                    onEnvChange={(env) => fetchConfigs(env)}
                />
                <div style={{ marginLeft: "auto", display: "flex", alignItems: "center", gap: 8 }}>
                    <Button type="primary" onClick={handleSave} disabled={modifiedValues.size === 0}>
                        Save ({modifiedValues.size})
                    </Button>
                </div>
            </div>

            <Row gutter={24} style={{ alignItems: "flex-start" }}>
                <Col span={20}>
                    <Input
                        placeholder="filter config"
                        prefix={<SearchOutlined />}
                        value={filter}
                        onChange={(e) => setFilter(e.target.value)}
                        style={{ marginBottom: 16 }}
                        allowClear
                    />

                    {Object.entries(filteredGroups).map(([group, items]) => (
                        <Card
                            key={group}
                            ref={(el) => (groupRefs.current[group] = el)}
                            id={group}
                            title={
                                <span
                                    style={{
                                        backgroundColor: highlighted === group ? "yellow" : "transparent",
                                        transition: "background-color 0.5s ease",
                                        padding: 4,
                                        borderRadius: 4,
                                    }}
                                >
                  {group}
                </span>
                            }
                            style={{ marginBottom: 24 }}
                        >
                            {items.map((cfg, idx) => (
                                <div key={cfg.key}>
                                    <ConfigRow
                                        title={cfg.key}
                                        description={cfg.usage}
                                        control={renderControl(cfg)}
                                        locked={!cfg.writable}
                                        query={filter}
                                        isChanged={modifiedValues.has(cfg.key)}
                                    />
                                    {idx < items.length - 1 && <Divider />}
                                </div>
                            ))}
                        </Card>
                    ))}
                </Col>

                <Col span={4}>
                    <Card title="Groups">
                        <div style={{ position: "sticky", top: 20 }}>
                            <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
                                {Object.keys(grouped).map((group) => (
                                    <a key={group} onClick={() => handleScrollToGroup(group)}>
                                        {group.toLowerCase()}
                                    </a>
                                ))}
                            </div>
                        </div>
                    </Card>
                </Col>
            </Row>
        </>
    );
};

export default ConfigsList;
