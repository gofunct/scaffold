package scaffold


var GrpcServerTemplate = `
package cmd


import (
	"github.com/gofunct/gotasks/runtime/grpc"
	"github.com/spf13/cobra"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "start a grpc server with config",
	Run: Serve(),
}

func init() { RootCmd.AddCommand(grpcCmd) }

func Serve() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", vi.VString("grpc_port"))
		if err != nil {
			log.Fatal("Failed to listen:"+vi.VString("grpc_port"), err)
		}
		tracer, closer, err := Trace(log.JZap)
		if err != nil {
			log.Fatal("Cannot initialize Jaeger Tracer %s", zap.Error(err))
		}
		defer closer.Close()

		// Set GRPC Interceptors
		server := NewServer(tracer)
		
		//Register your grpc service server here 
		// ex: api.RegisterTodoServiceServer(server, &db.Store{DB: NewDB()})

		mux := NewMux()
		log.Zap.Debug("Starting debug service..", zap.String("grpc_debug_port", vi.VString("grpc_debug_port")))
		go func() { http.ListenAndServe(vi.VString("grpc_debug_port"), mux) }()

		log.Zap.Debug("Starting grpc service..", zap.String("grpc_port", vi.VString("grpc_port")))
		server.Serve(lis)
	}
}

`