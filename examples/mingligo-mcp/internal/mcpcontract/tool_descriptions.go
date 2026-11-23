package mcpcontract

// toolDescriptions is the single source for MingliGo tool descriptions.
// Typed registration consumes this table and ToolNames derives the public
// surface from it.
var toolDescriptions = map[string]string{
	"time_now":           `按 IANA 时区返回服务端当前标准时间（无 DST）；空值、无效或未知 timezone 直接报错。所有「现在/今天/当前年份」运算以本返回值为唯一时间基准。用户出生时间不走本工具：民用/农历/时辰/日界模糊时先用 subject_candidates；地点或时区未知先用 city_info；黄历择日用 huangli_days。`,
	"huangli_days":       `黄历择日查日。用户询问某事项宜忌日期、某日吉凶或黄历宜忌事实时调用；从 standard_time 起返回连续 count 天，按 event 逐日附适配结论与判定依据，范围内无适配日也不报错。命理合参（八字/紫微本命、流年）不走本工具，改用对应排盘与查询工具。`,
	"city_info":          `根据城市名返回城市信息候选（内置表，不访问外链）：每项含名称、行政区、层级、经纬度、IANA 时区与可选经度范围；locale 只本地化名称展示，不改变匹配集合。长度 1 可直接使用；多于 1（同名区县等）必须由用户确认，不得静默取第一个。未命中时请询问附近大城市，或让用户直接提供经纬度和时区。`,
	"subject_candidates": `把民用时间事实归一化为 canonical subject 候选：接受公历/农历壁钟、时辰或同日区间；longitude_range 可选，表示地点边界不确定性。返回候选列表，input 候选恒定第一，边界派生候选按时间排序，每个候选附 reasons；longitude_range 省略时只用精确 longitude，非法时间、时区或范围直接报错。时间明确时直接组装下游 subject；模糊或触及时辰、日界、交节、夏令时边界时先调用本工具并等待确认。地点未知先用 city_info，不得静默选择候选或编造缺失时间/地点。`,
	"bazi_chart":         `排八字本命盘并返回确定性结构事实。subject 可由用户明确事实组装，也可取自用户确认后的 subject_candidates 候选；时间模糊或多候选时先归一化并确认。需要四柱、大运、十神、藏干、神煞、刑冲合会、扶抑/调候或格局证据时调用。只要原始本命盘，不用八字紫微合参断语；后者用 natal_query，紫微宫位细节用 ziwei_chart。`,
	"bazhai_chart":       `八宅排盘：由已校准 subject 推导出生年与性别，排命卦、四吉四凶方与出生年紫白飞星。subject 必须来自用户确认事实或 subject_candidates；时间或性别非法时返回错误，不猜测出生时辰。住宅风水评估、命卦方位吉凶类问题调用；门/主/灶配合改用 bazhai_layout。八宅不与八字四柱合参。`,
	"bazhai_layout":      `八宅门主灶配合：评估命卦与门/主/灶卦的游年九星、方向与吉凶。门/主/灶卦必须来自建筑方位或用户确认事实，本工具不反推方位。已由 bazhai_chart 得到命卦后调用；输出是八宅学派确定性参考判断，不是施工或安全建议；玄空飞星用 xuankong 工具。`,
	"xuankong_chart":     `玄空飞星起盘：按宅体事实日期（建成、入住或改宅）与二十四山坐向，排九宫飞星盘。评估阳宅固定风水格局时调用；流年飞星变化改用 xuankong_liunian。注意起盘日期不是命主出生时间。`,
	"xuankong_liunian":   `玄空流年飞星：在宅运盘上叠加指定年份，返回流年飞星盘；count=1（默认）单年，count>1 连续批量。评估阳宅某一年风水变化时调用；固定格局用 xuankong_chart，尚未起宅运盘的也先调用 xuankong_chart。`,
	"liuyao_qigua":       `六爻起卦事件创建。qi_gua_ji_lu 与 qi_gua_qing_jing（时间、经度、占问焦点）构成唯一不可变事件；起卦时按焦点确定性命理规则解析用神，并写入 qi_gua.qi_gua_qing_jing.qu_yong。qi_gua_ji_lu.fang_shi=auto 每次都是新占卜事件，tong_qian_ji_lu/yao_zhi 可复现既有卦。question 由调用方持有，不进入 qi_gua。返回完整 qi_gua，后续 liuyao_chart 必须原样消费。`,
	"liuyao_chart":       `六爻装卦。只消费 liuyao_qigua 返回的最小 qi_gua，读取 qi_gua.qi_gua_qing_jing.qu_yong，装卦并输出纳甲、六亲、六兽、旺衰、格局、应期与矛盾信号。本方法不创建起卦事件，也不输出解读。`,
	"bazi_bond":          `八字合盘返回确定性关系事实，不做合婚评级。a/b 必须来自用户确认事实或 subject_candidates 归一化候选；非法、不完整或未确认的模糊 subject 直接报错。用户询问合婚、合作、亲子等关系并需要八字侧证据时调用；需要紫微宫位主星对照用 ziwei_bond，两者常配合使用。`,
	"bazi_liupan":        `八字流盘：按 scope 返回 liu_nian/liu_yue/liu_ri/liu_shi/xiao_yun 对应领域数组，字段见 output schema。按粒度选择：逐年用 liu_nian，逐月用 liu_yue，具体日期用 liu_ri，具体时辰用 liu_shi，童年或小运趋势用 xiao_yun。liu_nian/liu_yue/liu_ri/liu_shi 的 count=1 单层、count>1 连续批量；xiao_yun 恒返回两个学派，count 只控制各学派柱数；liu_ri/liu_shi 必须带经度，非法 scope/count/time 或缺失必需经度直接报错。只返回确定性流盘事实，不做断语；本命盘用 bazi_chart，八字紫微逐年合参断语用 natal_liunian_query。`,
	"ziwei_chart":        `紫微斗数排本命盘并返回确定性宫位事实。subject 可由用户明确事实组装，也可取自用户确认后的 subject_candidates 候选；时间模糊或多候选时先归一化并确认。需要十二宫主星、亮度、四化、杂曜、长生博士十二神、小限或大限时调用。两人合盘用 ziwei_bond，流盘用 ziwei_liupan，八字紫微断语合参用 natal_query。`,
	"ziwei_liupan":       `紫微流盘：按 scope 返回 liu_nian/liu_yue/liu_ri/liu_shi 对应的领域数组，字段见 output schema。count=1（默认）单层，count>1 连续批量；standard_time 为当地标准时间，liu_shi 必须带经度。只返回确定性流盘事实，不做断语；本命盘用 ziwei_chart，八字紫微逐年合参断语用 natal_liunian_query。`,
	"ziwei_bond":         `紫微合盘原始事实。a/b 必须是用户确认事实或 subject_candidates 归一化后的主体，不得手造模糊时间。用户询问两人关系并需要紫微侧证据时调用；八字干支与十神视角改用 bazi_bond，两者常配合使用。并排返回双方命宫、夫妻宫、子女宫主星；不做跨盘投射和合婚评级。`,
	"qimen_chart":        `奇门排盘。method 显式选择时/刻/日/月/年家与转盘/飞盘/鸣法/金函玉镜的合法组合；时间、地点、用神保持直接参数。`,
}
