package main

import (
	"bytes"
	"strings"
	"text/template"
)

var asynqTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{.ServiceType}}JobServer interface {
{{- range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) (error)
{{- end}}
}

func Register{{.ServiceType}}JobServer(mux *asynq.ServeMux, srv {{.ServiceType}}JobServer) {
	{{- range .Methods}}
	mux.HandleFunc("{{.Typename}}", _{{$svrType}}_{{.Name}}_Job_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}_Job_Handler(srv {{$svrType}}JobServer) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		var in {{.Request}}
		if err := proto.Unmarshal(task.Payload(), &in); err != nil {
			return err
		}
		err := srv.{{.Name}}(ctx, &in)
		return err
	}
}
{{end}}

type {{.ServiceType}}SvcJob struct {}
var {{.ServiceType}}Job {{.ServiceType}}SvcJob

{{range .MethodSets}}
func (j *UserSvcJob) {{.Name}}(in *{{.Request}}, opts ...asynq.Option) (*asynq.Task, error) {
	payload, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask("{{.Typename}}", payload, opts...)
	return task, nil
}
{{end}}

type {{.ServiceType}}JobClient interface {
{{- range .MethodSets}}
	{{.Name}}(ctx context.Context, req *{{.Request}}, opts ...asynq.Option) (info *asynq.TaskInfo, err error) 
{{- end}}
}

type {{.ServiceType}}JobClientImpl struct{
	cc *asynq.Client
}
	
func New{{.ServiceType}}JobClient (client *asynq.Client) {{.ServiceType}}JobClient {
	return &{{.ServiceType}}JobClientImpl{client}
}

{{range .MethodSets}}
func (c *{{$svrType}}JobClientImpl) {{.Name}}(ctx context.Context, in *{{.Request}}, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	task, err := {{$svrType}}Job.{{.Name}}(in, opts...)
	if err != nil {
		return nil, err
	}
	info, err := c.cc.Enqueue(task)
	if err != nil {
		return nil, err
	}
	return info, nil
}
{{end}}
`

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc
}

type methodDesc struct {
	// method
	Name    string
	Num     int
	Request string
	Reply   string
	// asynq rule
	Typename string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("asynq").Parse(strings.TrimSpace(asynqTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
