package backups

import (
	"fmt"
	"net/http"

	"github.com/eduardolat/pgbackweb/internal/database/dbgen"
	"github.com/eduardolat/pgbackweb/internal/service/backups"
	"github.com/eduardolat/pgbackweb/internal/util/echoutil"
	"github.com/eduardolat/pgbackweb/internal/util/paginateutil"
	"github.com/eduardolat/pgbackweb/internal/util/timeutil"
	"github.com/eduardolat/pgbackweb/internal/validate"
	"github.com/eduardolat/pgbackweb/internal/view/web/component"
	"github.com/eduardolat/pgbackweb/internal/view/web/htmx"
	"github.com/labstack/echo/v4"
	"github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

func (h *handlers) listBackupsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var formData struct {
		Page int `query:"page" validate:"required,min=1"`
	}
	if err := c.Bind(&formData); err != nil {
		return htmx.RespondToastError(c, err.Error())
	}
	if err := validate.Struct(&formData); err != nil {
		return htmx.RespondToastError(c, err.Error())
	}

	pagination, backups, err := h.servs.BackupsService.PaginateBackups(
		ctx, backups.PaginateBackupsParams{
			Page:  formData.Page,
			Limit: 8,
		},
	)
	if err != nil {
		return htmx.RespondToastError(c, err.Error())
	}

	return echoutil.RenderGomponent(
		c, http.StatusOK, listBackups(pagination, backups),
	)
}

func listBackups(
	pagination paginateutil.PaginateResponse,
	backups []dbgen.BackupsServicePaginateBackupsRow,
) gomponents.Node {
	yesNoSpan := func(b bool) gomponents.Node {
		if b {
			return component.SpanText("Yes")
		}
		return component.SpanText("No")
	}

	trs := []gomponents.Node{}
	for _, backup := range backups {
		trs = append(trs, html.Tr(
			html.Td(
				html.Class("w-[40px]"),
				html.Div(
					html.Class("flex justify-start space-x-1"),
					editBackupButton(backup),
					deleteBackupButton(backup.ID),
				),
			),
			html.Td(
				html.Div(
					html.Class("flex items-center space-x-2"),
					html.Div(
						html.Class("tooltip tooltip-right"),
						gomponents.If(backup.IsActive, html.Data("tip", "Active")),
						gomponents.If(!backup.IsActive, html.Data("tip", "Inactive")),
						component.EnabledPing(backup.IsActive),
					),
					component.SpanText(backup.Name),
				),
			),
			html.Td(component.SpanText(backup.DatabaseName)),
			html.Td(component.SpanText(backup.DestinationName)),
			html.Td(
				html.Class("font-mono"),
				html.Div(
					html.Class("flex flex-col items-start space-y-1"),
					component.SpanText(backup.CronExpression),
					component.SpanText(backup.TimeZone),
				),
			),
			html.Td(component.SpanText(fmt.Sprintf("%d days", backup.RetentionDays))),
			html.Td(yesNoSpan(backup.OptDataOnly)),
			html.Td(yesNoSpan(backup.OptSchemaOnly)),
			html.Td(yesNoSpan(backup.OptClean)),
			html.Td(yesNoSpan(backup.OptIfExists)),
			html.Td(yesNoSpan(backup.OptCreate)),
			html.Td(yesNoSpan(backup.OptNoComments)),
			html.Td(component.SpanText(
				backup.CreatedAt.Format(timeutil.LayoutYYYYMMDDHHMMSSPretty),
			)),
		))
	}

	if pagination.HasNextPage {
		trs = append(trs, html.Tr(
			htmx.HxGet(fmt.Sprintf(
				"/dashboard/backups/list?page=%d", pagination.NextPage,
			)),
			htmx.HxTrigger("intersect once"),
			htmx.HxSwap("afterend"),
			htmx.HxIndicator("#list-backups-loading"),
		))
	}

	return component.RenderableGroup(trs)
}