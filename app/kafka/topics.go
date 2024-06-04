package kafka

const (
	NOTIFICATION_TOPIC = "notification"
	NOTIFICATION_GROUP = "notification-group"
	EVALUATION         = "evaluation"
	USER_MUTATION      = "user-mutation"
)

/*
	export const EventValidator = z.object({
	  message: z.string(),
	  enableEmail: z.boolean().default(false),
	  enablePush: z.boolean().default(true),
	  title: z.string(),
	  badge: z.string().optional().nullable(),
	  userId: z.string().optional().nullable(),
	  role: z.string(),
	  group: z.string(),
	  pushTo: z.nativeEnum(PushTo),
	  year: z.string(),
	});
*/
const (
	USER_NOTIFICATION_TYPE  = "user"
	GROUP_NOTIFICATION_TYPE = "group"
	PROMO_NOTIFICATION_TYPE = "promo"
	ALL_NOTIFICATION_TYPE   = "all"
)

type NotificationEvent struct {
	Message     string  `json:"message"`
	EnableEmail bool    `json:"enableEmail"`
	EnablePush  bool    `json:"enablePush"`
	Title       string  `json:"title"`
	Badge       *string `json:"badge"`
	UserId      string  `json:"userId"`
	Role        string  `json:"role"`
	Group       string  `json:"group"`
	PushTo      string  `json:"pushTo"`
	Year        string  `json:"year"`
}
