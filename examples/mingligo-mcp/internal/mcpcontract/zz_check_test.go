package mcpcontract
import ("encoding/json";"os";"testing";"time";"github.com/ml8s/mingligo/pkg/bazi";"github.com/ml8s/mingligo/pkg/tianwen")
func TestCheckFlatFields(t *testing.T){
  st,_:=time.Parse(time.RFC3339,"1990-05-20T11:49:00+08:00")
  c,err:=bazi.ComputeChartAt(st,st); if err!=nil{t.Fatal(err)}
  _=tianwen.JieQi(0)
  b,_:=json.Marshal(c)
  var m map[string]any; json.Unmarshal(b,&m)
  out,_:=json.Marshal(map[string]any{"spouse_palace_state":m["spouse_palace_state"],"day_branch_type":m["day_branch_type"],"has_derived":m["derived"]!=nil})
  os.WriteFile("/tmp/opencode/flat_fields.txt",out,0644)
}
