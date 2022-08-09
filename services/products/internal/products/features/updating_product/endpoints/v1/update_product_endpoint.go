package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/meysamhadeli/shop-golang-microservices/pkg/mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/tracing"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/delivery"
	"github.com/meysamhadeli/shop-golang-microservices/services/products/internal/products/features/updating_product"
)

type updateProductEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewUpdateProductEndpoint(productEndpointBase *delivery.ProductEndpointBase) *updateProductEndpoint {
	return &updateProductEndpoint{productEndpointBase}
}

func (ep *updateProductEndpoint) MapRoute() {
	ep.ProductsGroup.PUT("/:id", ep.updateProduct())
}

// UpdateProduct
// @Tags Products
// @Summary Update product
// @Description Update existing product
// @Accept json
// @Produce json
// @Param UpdateProductRequestDto body updating_product.UpdateProductRequestDto true "Product data"
// @Param id path string true "Product ID"
// @Success 204
// @Router /api/v1/products/{id} [put]
func (ep *updateProductEndpoint) updateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.Metrics.UpdateProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "updateProductEndpoint.updateProduct")
		defer span.Finish()

		request := &updating_product.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Log.Warn("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		command := updating_product.NewUpdateProduct(request.ProductID, request.Name, request.Description, request.Price)

		if err := ep.Validator.StructCtx(ctx, command); err != nil {
			ep.Log.Warn("validate", err)
			tracing.TraceErr(span, err)
			return err
		}

		_, err := mediatr.Send[*mediatr.Unit](ctx, command)

		if err != nil {
			ep.Log.Warnf("UpdateProduct", err)
			tracing.TraceErr(span, err)
			return err
		}

		ep.Log.Infof("(product updated) id: {%s}", request.ProductID)

		return c.NoContent(http.StatusNoContent)
	}
}
