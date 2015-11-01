package lib

/*
import (
	pb "github.com/clawio/service.auth/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.AuthClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c Client) Authenticate(ctx context.Context, username, password string) (string, error) {
	areq := &pb.AuthRequest{}
	areq.Username = username
	areq.Password = password

	res, err := c.client.Authenticate(ctx, areq)
	if err != nil {
		return "", err
	}
	return res.Token, nil
}

func (c Client) Close() error {
	return c.conn.Close()
}
*/
