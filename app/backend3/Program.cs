var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/rolldice", () => {
    return new Random().Next( 1, 6 );
});

app.Run();
