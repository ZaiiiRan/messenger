package httputil

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var marshaler = protojson.MarshalOptions{
	UseProtoNames:   false,
	EmitUnpopulated: true,
}

var unmarshaler = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

// Bind decodes the JSON request body into the given protobuf message. An empty
// body is treated as an empty message (handy for parameterless RPCs).
func Bind(c *fiber.Ctx, msg proto.Message) error {
	body := c.Body()
	if len(body) == 0 {
		return nil
	}
	if err := unmarshaler.Unmarshal(body, msg); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	return nil
}

// Proto writes a protobuf message as a JSON response with status 200.
func Proto(c *fiber.Ctx, msg proto.Message) error {
	return ProtoStatus(c, fiber.StatusOK, msg)
}

// ProtoStatus writes a protobuf message as a JSON response with a custom status.
func ProtoStatus(c *fiber.Ctx, status int, msg proto.Message) error {
	data, err := marshaler.Marshal(msg)
	if err != nil {
		return err
	}
	c.Status(status)
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Send(data)
}

// RawProto marshals a protobuf message into raw JSON so it can be embedded into
// a composite response payload.
func RawProto(msg proto.Message) (json.RawMessage, error) {
	if msg == nil {
		return json.RawMessage("null"), nil
	}
	data, err := marshaler.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}

// JSON writes an arbitrary value as a JSON response with status 200.
func JSON(c *fiber.Ctx, payload any) error {
	return c.Status(fiber.StatusOK).JSON(payload)
}
